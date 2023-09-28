package domain

import (
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

const InvalidPlaceId = MaxPlaceCount
const InvalidParticipantId = MaxParticipantCount
const InvalidMessageId = MaxMessageCount

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

func InvalidTransition() Transition {
	return Transition{
		IsValid:        false,
		IncomingPlaces: [MaxBranchingFactor]PlaceId{InvalidPlaceId, InvalidPlaceId},
		OutgoingPlaces: [MaxBranchingFactor]PlaceId{InvalidPlaceId, InvalidPlaceId},
		Participant:    InvalidParticipantId,
		Message:        InvalidMessageId,
	}
}

type ModelId = string

type Model struct {
	Id               ModelId
	Hash             Hash
	Name             string
	PlaceCount       uint8
	ParticipantCount uint8
	MessageCount     uint8
	StartPlaces      [MaxStartPlaceCount]PlaceId
	EndPlaces        [MaxEndPlaceCount]PlaceId
	Transitions      [MaxTransitionCount]Transition
}

func (model *Model) Instantiate(publicKeys []PublicKey) (Instance, error) {
	if int(model.ParticipantCount) != len(publicKeys) {
		return Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the model %s", model.Id)
	}
	var tokenCounts [MaxPlaceCount]int8
	for i := model.PlaceCount; i < MaxPlaceCount; i++ {
		tokenCounts[i] = InvalidTokenCount
	}
	for _, startPlace := range model.StartPlaces {
		if startPlace != InvalidPlaceId {
			tokenCounts[startPlace] = 1
		}
	}
	var messageHashes [MaxMessageCount]Hash
	for i := model.MessageCount; i < MaxMessageCount; i++ {
		messageHashes[i] = InvalidHash()
	}
	var publicKeysFixedSize [MaxParticipantCount]PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := model.ParticipantCount; i < MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = InvalidPublicKey()
	}
	instance := Instance{
		Model:         model.Id,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeysFixedSize,
		MessageHashes: messageHashes,
		UpdatedAt:     time.Now().Unix(),
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
	return Transition{}, fmt.Errorf("transition %s not found in model %s", id, model.Id)
}
