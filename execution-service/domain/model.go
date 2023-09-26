package domain

import "fmt"

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

const DefaultPlaceId = MaxPlaceCount
const DefaultParticipantId = MaxParticipantCount
const DefaultMessageId = MaxMessageCount

type Transition struct {
	Id             string
	Name           string
	IsInitialized  bool
	IncomingPlaces [MaxBranchingFactor]PlaceId
	OutgoingPlaces [MaxBranchingFactor]PlaceId
	Participant    ParticipantId
	Message        MessageId
}

var DefaultTransition = Transition{
	IsInitialized:  false,
	IncomingPlaces: [MaxBranchingFactor]PlaceId{DefaultPlaceId, DefaultPlaceId},
	OutgoingPlaces: [MaxBranchingFactor]PlaceId{DefaultPlaceId, DefaultPlaceId},
	Participant:    DefaultParticipantId,
	Message:        DefaultMessageId,
}

type Model struct {
	Id               string
	Hash             []byte
	PlaceCount       uint8
	ParticipantCount uint8
	MessageCount     uint8
	StartPlaces      [MaxStartPlaceCount]PlaceId
	EndPlaces        [MaxEndPlaceCount]PlaceId
	Transitions      [MaxTransitionCount]Transition
	Salt             []byte
}

func (model *Model) Instantiate(publicKeys []PublicKey) (Instance, error) {
	if int(model.ParticipantCount) != len(publicKeys) {
		return Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the model %s", model.Id)
	}
	var tokenCounts [MaxPlaceCount]int8
	for _, startPlace := range model.StartPlaces {
		if startPlace != DefaultPlaceId {
			tokenCounts[startPlace] = 1
		}
	}
	var messageHashes [MaxMessageCount]MessageHash
	for i := 0; i < MaxMessageCount; i++ {
		messageHashes[i] = DefaultMessageHash
	}
	var publicKeysFixedSize [MaxParticipantCount]PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := model.ParticipantCount; i < MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = DefaultPublicKey
	}
	instance := Instance{
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeysFixedSize,
		MessageHashes: messageHashes,
	}
	instance.ComputeHash()
	return instance, nil
}
