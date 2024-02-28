package domain

import (
	"fmt"
	"math"
	"time"
)

const BasePlaceCount = 16
const BaseParticipantCount = 16
const BaseMessageCount = 16
const BaseStartPlaceCount = 1
const BaseEndPlaceCount = 2
const BaseTransitionCount = 16
const BaseBranchingFactor = 4
const BaseConstraintMessageCount = 4

const factor = 4

const MaxPlaceCount = BasePlaceCount * factor
const MaxParticipantCount = BaseParticipantCount * factor
const MaxMessageCount = BaseMessageCount * factor
const MaxStartPlaceCount = BaseStartPlaceCount * factor
const MaxEndPlaceCount = BaseEndPlaceCount * factor
const MaxTransitionCount = BaseTransitionCount * factor
const MaxBranchingFactor = BaseBranchingFactor
const MaxConstraintMessageCount = BaseConstraintMessageCount

var MaxParticipantDepth = int(math.Log2(MaxParticipantCount))
var MaxTransitionDepth = int(math.Log2(MaxTransitionCount))
var MaxEndPlaceDepth = int(math.Log2(MaxEndPlaceCount))

type PlaceId = uint16
type ParticipantId = uint16
type ModelMessageId = uint16
type IntegerType = int32
type ComparisonOperator = uint8
type ModelId = string
type TransitionId = string

const (
	OperatorEqual              = 0
	OperatorGreaterThan        = 1
	OperatorLessThan           = 2
	OperatorGreaterThanOrEqual = 3
	OperatorLessThanOrEqual    = 4
)

var ValidComparisonOperators = []ComparisonOperator{OperatorEqual, OperatorGreaterThan, OperatorLessThan, OperatorGreaterThanOrEqual, OperatorLessThanOrEqual}

const OutOfBoundsPlaceId = PlaceId(MaxPlaceCount)
const EmptyParticipantId = ParticipantId(MaxParticipantCount)
const EmptyMessageId = ModelMessageId(MaxMessageCount)

type Model struct {
	Hash             SaltedHash
	PlaceCount       uint16
	ParticipantCount uint16
	MessageCount     uint16
	StartPlaces      []PlaceId
	EndPlaces        []PlaceId
	Transitions      []Transition
	CreatedAt        int64
}

type Transition struct {
	Id                    TransitionId
	Name                  string
	IncomingPlaces        []PlaceId
	OutgoingPlaces        []PlaceId
	InitiatingParticipant ParticipantId
	RespondingParticipant ParticipantId
	InitiatingMessage     ModelMessageId
	RespondingMessage     ModelMessageId
	Constraint            Constraint
}

// ax + by + c = 0
type Constraint struct {
	Coefficients       []IntegerType
	MessageIds         []ModelMessageId
	Offset             IntegerType
	ComparisonOperator ComparisonOperator
}

type ConstraintInput struct {
	Messages []Message
}

func OutOfBoundsTransition() Transition {
	return Transition{
		IncomingPlaces:        make([]PlaceId, 0),
		OutgoingPlaces:        make([]PlaceId, 0),
		InitiatingParticipant: EmptyParticipantId,
		RespondingParticipant: EmptyParticipantId,
		InitiatingMessage:     EmptyMessageId,
		Constraint:            EmptyConstraint(),
	}
}

func EmptyConstraint() Constraint {
	return Constraint{
		Coefficients: make([]IntegerType, 0),
		MessageIds:   make([]ModelMessageId, 0),
	}
}

func EmptyConstraintInput() ConstraintInput {
	return ConstraintInput{
		Messages: make([]Message, 0),
	}
}

func (model *Model) Id() ModelId {
	return model.Hash.String()
}

func (model *Model) Instantiate(publicKeys []PublicKey) (Instance, error) {
	if len(publicKeys) > MaxParticipantCount {
		return Instance{}, fmt.Errorf("there are too many public keys")
	}
	if int(model.ParticipantCount) != len(publicKeys) {
		return Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the model %s", model.Id())
	}
	tokenCounts := make([]int8, model.PlaceCount)
	for _, startPlace := range model.StartPlaces {
		tokenCounts[startPlace] = 1
	}
	messageHashes := make([]Hash, model.MessageCount)
	for i := 0; i < int(model.MessageCount); i++ {
		messageHashes[i] = EmptyHash()
	}
	instance := Instance{
		Model:         model.Hash.Hash,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
		CreatedAt:     time.Now().Unix(),
	}
	instance.UpdateHash()
	return instance, nil
}

func (model *Model) FindTransitionById(id TransitionId) (Transition, error) {
	for _, transition := range model.Transitions {
		if transition.Id == id {
			return transition, nil
		}
	}
	return Transition{}, fmt.Errorf("transition %s not found in model %s", id, model.Id())
}

func (model *Model) FindNextParticipants(transition Transition) []ParticipantId {
	participants := make([]ParticipantId, 0)
	for _, nextTransition := range model.Transitions {
		intersection := intersect(transition.OutgoingPlaces, nextTransition.IncomingPlaces)
		if len(intersection) > 0 && nextTransition.InitiatingParticipant != EmptyParticipantId {
			participants = append(participants, nextTransition.InitiatingParticipant)
		}
	}
	return participants
}

func intersect(set1 []PlaceId, set2 []PlaceId) []PlaceId {
	result := make([]PlaceId, 0)
	hash := make(map[PlaceId]bool)
	for _, v := range set1 {
		hash[v] = true
	}
	for _, v := range set2 {
		if hash[v] {
			result = append(result, v)
			hash[v] = false
		}
	}
	return result
}
