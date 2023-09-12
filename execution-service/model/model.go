package model

type PlaceId uint

type TransitionId string

type Transition struct {
	Id             TransitionId
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type ModelId string

type Model struct {
	Id          ModelId
	PlaceCount  uint
	StartPlace  PlaceId
	Transitions []Transition
}
