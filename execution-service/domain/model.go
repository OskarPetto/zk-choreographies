package domain

import "fmt"

const MaxPlaceCount = 64
const MaxParticipantCount = 8
const MaxMessageCount = 64
const MaxStartPlaceCount = 4
const MaxEndPlaceCount = 4
const MaxTransitionCount = 64
const MaxBranchingFactor = 4

type PlaceId = uint
type ParticipantId = uint
type MessageId = uint

type Transition struct {
	Id             string
	Name           string
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
	Participant    ParticipantId
	Message        MessageId
}

type Model struct {
	Id               string
	PlaceCount       uint
	ParticipantCount uint
	MessageCount     uint
	StartPlaces      []PlaceId
	EndPlaces        []PlaceId
	Transitions      []Transition
}

func (model *Model) Instantiate(publicKeys []PublicKey) (Instance, error) {
	if int(model.ParticipantCount) != len(publicKeys) {
		return Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the model %s", model.Id)
	}
	tokenCounts := make([]int, model.PlaceCount)
	for _, startPlace := range model.StartPlaces {
		tokenCounts[startPlace] = 1
	}
	messageHashes := make([]MessageHash, model.MessageCount)
	instance := Instance{
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
	}
	instance.ComputeHash()
	return instance, nil
}
