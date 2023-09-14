package workflow

type PlaceId = uint

type Transition struct {
	Id             string
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type PetriNet struct {
	Id          string
	PlaceCount  uint
	StartPlace  PlaceId
	Transitions []Transition
}

type Instance struct {
	Id          string
	PetriNet    string
	TokenCounts []int
}
