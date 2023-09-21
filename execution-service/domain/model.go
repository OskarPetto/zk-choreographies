package domain

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

type PetriNet struct {
	Id               string
	StartPlace       PlaceId
	EndPlace         PlaceId
	PlaceCount       uint
	ParticipantCount uint
	Transitions      []Transition
}
