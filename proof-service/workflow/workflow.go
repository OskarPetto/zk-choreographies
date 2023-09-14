package workflow

type PlaceId = uint8

type Instance struct {
	Id          string
	PetriNet    string
	TokenCounts []int8
}

type Transition struct {
	Id             string
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type PetriNet struct {
	Id          string
	PlaceCount  uint8
	StartPlace  PlaceId
	Transitions []Transition
}
