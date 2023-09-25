package domain

import "fmt"

const MaxPlaceCount = 100
const MaxParticipantCount = 20
const MaxMessageCount = 80
const MaxTransitionCount = 100
const MaxBranchingFactor = 3
const MaxEndPlaceCount = 5

type PlaceId = uint
type ParticipantId = uint
type MessageId = uint

type Transition struct {
	Id                 string
	Name               string
	IncomingPlaces     []PlaceId
	OutgoingPlaces     []PlaceId
	ParticipantIsValid bool
	Participant        ParticipantId
	MessageIsValid     bool
	Message            MessageId
}

type Model struct {
	Id               string
	PlaceCount       uint
	ParticipantCount uint
	MessageCount     uint
	StartPlace       PlaceId
	EndPlaces        []PlaceId
	Transitions      []Transition
}

func (model *Model) Instantiate(publicKeys []PublicKey) (Instance, error) {
	if int(model.ParticipantCount) != len(publicKeys) {
		return Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the model %s", model.Id)
	}
	tokenCounts := make([]int, model.PlaceCount)
	tokenCounts[model.StartPlace] = 1
	messageHashes := make([]MessageHash, model.MessageCount)
	instance := Instance{
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
	}
	instance.ComputeHash()
	return instance, nil
}
