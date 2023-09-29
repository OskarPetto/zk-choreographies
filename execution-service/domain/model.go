package domain

import (
	"execution-service/utils"
	"fmt"
	"time"
)

const MaxPlaceCount = 64
const MaxParticipantCount = 8
const MaxMessageCount = 64
const MaxStartPlaceCount = 1
const MaxEndPlaceCount = 2
const MaxTransitionCount = 64
const MaxBranchingFactor = 2

type PlaceId = uint8
type ParticipantId = uint8
type MessageId = uint8

const OutOfBoundsPlaceId = MaxPlaceCount
const EmptyParticipantId = MaxParticipantCount
const EmptyMessageId = MaxMessageCount

type TransitionId = string

type Transition struct {
	Id             TransitionId
	Name           string
	IsValid        bool
	IncomingPlaces [MaxBranchingFactor]PlaceId
	OutgoingPlaces [MaxBranchingFactor]PlaceId
	Participant    ParticipantId
	Message        MessageId
}

func OutOfBoundsTransition() Transition {
	return Transition{
		IsValid:        false,
		IncomingPlaces: [MaxBranchingFactor]PlaceId{OutOfBoundsPlaceId, OutOfBoundsPlaceId},
		OutgoingPlaces: [MaxBranchingFactor]PlaceId{OutOfBoundsPlaceId, OutOfBoundsPlaceId},
		Participant:    EmptyParticipantId,
		Message:        EmptyMessageId,
	}
}

type ModelId = string

type Model struct {
	Hash             Hash
	Choreography     string
	PlaceCount       uint8
	ParticipantCount uint8
	MessageCount     uint8
	StartPlaces      [MaxStartPlaceCount]PlaceId
	EndPlaces        [MaxEndPlaceCount]PlaceId
	Transitions      [MaxTransitionCount]Transition
	CreatedAt        int64
}

func (model *Model) Id() ModelId {
	return utils.BytesToString(model.Hash.Value[:])
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
	var messageHashes [MaxMessageCount]Hash
	for i := 0; i < int(model.MessageCount); i++ {
		messageHashes[i] = EmptyHash()
	}
	for i := model.MessageCount; i < MaxMessageCount; i++ {
		messageHashes[i] = OutOfBoundsHash()
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
