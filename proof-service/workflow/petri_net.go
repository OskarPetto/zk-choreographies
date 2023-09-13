package workflow

type PlaceId uint

type TransitionId string

type Transition struct {
	Id             TransitionId
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type PetriNetId string

type PetriNet struct {
	Id          PetriNetId
	PlaceCount  uint
	StartPlace  PlaceId
	Transitions []Transition
}
