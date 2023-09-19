package workflow

const MaxPlaceCount = 100
const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type PlaceId = uint

type Transition struct {
	Id             string
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type PetriNet struct {
	Id          string
	StartPlace  PlaceId
	EndPlace    PlaceId
	PlaceCount  uint
	Transitions []Transition
}

type Instance struct {
	Id          string
	TokenCounts []int
}
