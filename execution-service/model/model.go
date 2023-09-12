package model

type PlaceId uint

type TransitionId string

type Transition struct {
	Id             TransitionId
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type Model struct {
	Id          string
	PlaceCount  uint
	StartPlace  PlaceId
	Transitions []Transition
}
