package domain

import (
	"fmt"
	"math"
	"time"
)

const MaxPlaceCount = 64
const MaxParticipantCount = 16
const MaxMessageCount = 64
const MaxStartPlaceCount = 1
const MaxEndPlaceCount = 16
const MaxTransitionCount = 64
const MaxBranchingFactor = 4

var MaxParticipantDepth = int(math.Log2(MaxParticipantCount))
var MaxTransitionDepth = int(math.Log2(MaxTransitionCount))
var MaxEndPlaceDepth = int(math.Log2(MaxEndPlaceCount))

type PlaceId = uint8
type ParticipantId = uint8
type MessageId = uint8
type VariableId = uint8

const OutOfBoundsPlaceId = PlaceId(MaxPlaceCount)
const EmptyParticipantId = ParticipantId(MaxParticipantCount)
const EmptyMessageId = MessageId(MaxMessageCount)

type TransitionId = string

type Transition struct {
	Id             TransitionId
	Name           string
	IncomingPlaces [MaxBranchingFactor]PlaceId
	OutgoingPlaces [MaxBranchingFactor]PlaceId
	Participant    ParticipantId
	Message        MessageId
	Constraint     Constraint
}

func OutOfBoundsTransition() Transition {
	var incomingPlaces [MaxBranchingFactor]PlaceId
	for i, _ := range incomingPlaces {
		incomingPlaces[i] = OutOfBoundsPlaceId
	}
	var outgoingPlaces [MaxBranchingFactor]PlaceId
	for i, _ := range outgoingPlaces {
		outgoingPlaces[i] = OutOfBoundsPlaceId
	}
	return Transition{
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
		Participant:    EmptyParticipantId,
		Message:        EmptyMessageId,
		Constraint:     EmptyConstraint(),
	}
}

func IsOutOfBoundsTransition(transition Transition) bool {
	outOfBoundsTransition := OutOfBoundsTransition()
	for i, incomingPlace := range transition.IncomingPlaces {
		if outOfBoundsTransition.IncomingPlaces[i] != incomingPlace {
			return false
		}
	}
	for i, outgoingPlace := range transition.OutgoingPlaces {
		if outOfBoundsTransition.OutgoingPlaces[i] != outgoingPlace {
			return false
		}
	}
	if transition.Participant != outOfBoundsTransition.Participant {
		return false
	}
	if transition.Message != outOfBoundsTransition.Message {
		return false
	}
	return IsEmptyConstraint(transition.Constraint)
}

type ModelId = string

type Model struct {
	Hash             Hash
	Choreography     string
	PlaceCount       uint8
	ParticipantCount uint8
	MessageCount     uint8
	VariableCount    uint8
	StartPlaces      [MaxStartPlaceCount]PlaceId
	EndPlaces        [MaxEndPlaceCount]PlaceId
	Transitions      [MaxTransitionCount]Transition
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
	var tokenCounts [MaxPlaceCount]int8
	for i := model.PlaceCount; i < MaxPlaceCount; i++ {
		tokenCounts[i] = OutOfBoundsTokenCount
	}
	for _, startPlace := range model.StartPlaces {
		if startPlace != OutOfBoundsPlaceId {
			tokenCounts[startPlace] = 1
		}
	}
	var messageHashes [MaxMessageCount][HashSize]byte
	for i := 0; i < int(model.MessageCount); i++ {
		messageHashes[i] = EmptyHash().Value
	}
	for i := model.MessageCount; i < MaxMessageCount; i++ {
		messageHashes[i] = OutOfBoundsHash().Value
	}
	var publicKeysFixedSize [MaxParticipantCount]PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := model.ParticipantCount; i < MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = OutOfBoundsPublicKey()
	}
	instance := Instance{
		Model:         model.Id(),
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeysFixedSize,
		MessageHashes: messageHashes,
		CreatedAt:     time.Now().Unix(),
	}
	instance.ComputeHash()
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
