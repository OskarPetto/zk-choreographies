package domain

import "fmt"

const MaxPlaceCount = 100
const MaxTransitionCount = 100
const MaxBranchingFactor = 3
const MaxParticipantCount = 20

type PlaceId = uint
type ParticipantId = uint

type Transition struct {
	Id                           string
	IsExecutableByAnyParticipant bool
	Participant                  ParticipantId
	IncomingPlaces               []PlaceId
	OutgoingPlaces               []PlaceId
}

type Model struct {
	Id               string
	StartPlace       PlaceId
	EndPlace         PlaceId
	PlaceCount       uint
	ParticipantCount uint
	Transitions      []Transition
}

func (model *Model) Instantiate(publicKeys [][]byte) (Instance, error) {
	if int(model.ParticipantCount) != len(publicKeys) {
		return Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the model %s", model.Id)
	}
	tokenCounts := make([]int, model.PlaceCount)
	tokenCounts[model.StartPlace] = 1
	instance := Instance{
		TokenCounts: tokenCounts,
		PublicKeys:  publicKeys,
	}
	instance.ComputeHash()
	return instance, nil
}
