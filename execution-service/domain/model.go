package domain

import (
	"fmt"
	"math"
	"time"
)

const MaxPlaceCount = 128
const MaxMessageCount = 128
const MaxTransitionCount = 256
const MaxParticipantCount = 32
const MaxStartPlaceCount = 2
const MaxEndPlaceCount = 16
const MaxBranchingFactor = 4

var MaxParticipantDepth = int(math.Log2(MaxParticipantCount))
var MaxTransitionDepth = int(math.Log2(MaxTransitionCount))
var MaxEndPlaceDepth = int(math.Log2(MaxEndPlaceCount))

type PlaceId = uint16
type ParticipantId = uint16
type ModelMessageId = uint16

const OutOfBoundsPlaceId = PlaceId(MaxPlaceCount)
const EmptyParticipantId = ParticipantId(MaxParticipantCount)
const EmptyMessageId = ModelMessageId(MaxMessageCount)

type TransitionId = string

type Transition struct {
	Id                    TransitionId
	Name                  string
	IncomingPlaces        []PlaceId
	OutgoingPlaces        []PlaceId
	InitiatingParticipant ParticipantId
	RespondingParticipant ParticipantId
	Message               ModelMessageId
	Constraint            Constraint
}

func OutOfBoundsTransition() Transition {
	return Transition{
		IncomingPlaces:        make([]PlaceId, 0),
		OutgoingPlaces:        make([]PlaceId, 0),
		InitiatingParticipant: EmptyParticipantId,
		RespondingParticipant: EmptyParticipantId,
		Message:               EmptyMessageId,
		Constraint:            EmptyConstraint(),
	}
}

type ModelId = string

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
