package domain

import (
	"bytes"
	"fmt"
	"math"
	"time"
)

const MaxPlaceCount = 256
const MaxParticipantCount = 256
const MaxMessageCount = 256
const MaxStartPlaceCount = 4
const MaxEndPlaceCount = 16
const MaxTransitionCount = 256
const MaxBranchingFactor = 8
const MaxConstraintMessageCount = 4

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
	Source           string
	PlaceCount       uint16
	ParticipantCount uint16
	MessageCount     uint16
	StartPlaces      []PlaceId
	EndPlaces        []PlaceId
	Transitions      []Transition
	CreatedAt        int64
}

type Transition struct {
	Id             TransitionId
	Name           string
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
	Sender         ParticipantId
	Recipient      ParticipantId
	Message        ModelMessageId
	Constraint     Constraint
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
		IncomingPlaces: make([]PlaceId, 0),
		OutgoingPlaces: make([]PlaceId, 0),
		Sender:         EmptyParticipantId,
		Recipient:      EmptyParticipantId,
		Message:        EmptyMessageId,
		Constraint:     EmptyConstraint(),
	}
}

func EmptyConstraint() Constraint {
	return Constraint{}
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
		if len(intersection) > 0 && nextTransition.Sender != EmptyParticipantId {
			participants = append(participants, nextTransition.Sender)
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

func (instance *Instance) EvaluateConstraint(constraint Constraint, input ConstraintInput) bool {
	if len(constraint.MessageIds) != len(input.Messages) {
		return false
	}
	lhs := constraint.Offset
	for i, message := range input.Messages {
		hash := message.Hash.Hash
		messageId := EmptyMessageId
		for i, messageHash := range instance.MessageHashes {
			if bytes.Equal(hash.Value[:], messageHash.Value[:]) {
				messageId = ModelMessageId(i)
				break
			}
		}
		if constraint.Coefficients[i] != 0 && messageId != constraint.MessageIds[i] {
			return false
		}
		lhs += constraint.Coefficients[i] * input.Messages[i].IntegerMessage
	}

	switch constraint.ComparisonOperator {
	case OperatorEqual:
		return lhs == 0
	case OperatorGreaterThan:
		return lhs > 0
	case OperatorLessThan:
		return lhs < 0
	case OperatorGreaterThanOrEqual:
		return lhs >= 0
	case OperatorLessThanOrEqual:
		return lhs <= 0
	}
	return false
}
