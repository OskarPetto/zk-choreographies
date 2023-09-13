package testdata

import (
	"execution-service/instance"
	"execution-service/petriNet"
)

func GetPetriNet1() petriNet.PetriNet {
	return petriNet.PetriNet{
		Id:         "conformance_example",
		StartPlace: 9,
		PlaceCount: 10,
		Transitions: []petriNet.Transition{
			petriNet.Transition{
				Id:             "As",
				IncomingPlaces: []petriNet.PlaceId{9},
				OutgoingPlaces: []petriNet.PlaceId{0},
			},
			petriNet.Transition{
				Id:             "Da1",
				IncomingPlaces: []petriNet.PlaceId{0},
				OutgoingPlaces: []petriNet.PlaceId{8},
			},
			petriNet.Transition{
				Id:             "Aa",
				IncomingPlaces: []petriNet.PlaceId{0},
				OutgoingPlaces: []petriNet.PlaceId{1, 2},
			},
			petriNet.Transition{
				Id:             "Fa",
				IncomingPlaces: []petriNet.PlaceId{1},
				OutgoingPlaces: []petriNet.PlaceId{5},
			},
			petriNet.Transition{
				Id:             "Sso",
				IncomingPlaces: []petriNet.PlaceId{2},
				OutgoingPlaces: []petriNet.PlaceId{3},
			},
			petriNet.Transition{
				Id:             "Ro",
				IncomingPlaces: []petriNet.PlaceId{3},
				OutgoingPlaces: []petriNet.PlaceId{4},
			},
			petriNet.Transition{
				Id:             "Co",
				IncomingPlaces: []petriNet.PlaceId{4},
				OutgoingPlaces: []petriNet.PlaceId{2},
			},
			petriNet.Transition{
				Id:             "Ao",
				IncomingPlaces: []petriNet.PlaceId{5, 4},
				OutgoingPlaces: []petriNet.PlaceId{6},
			},
			petriNet.Transition{
				Id:             "Do",
				IncomingPlaces: []petriNet.PlaceId{5, 4},
				OutgoingPlaces: []petriNet.PlaceId{7},
			},
			petriNet.Transition{
				Id:             "Aaa",
				IncomingPlaces: []petriNet.PlaceId{6},
				OutgoingPlaces: []petriNet.PlaceId{8},
			},
			petriNet.Transition{
				Id:             "Da2",
				IncomingPlaces: []petriNet.PlaceId{7},
				OutgoingPlaces: []petriNet.PlaceId{8},
			},
			petriNet.Transition{
				Id:             "Af",
				IncomingPlaces: []petriNet.PlaceId{8},
				OutgoingPlaces: []petriNet.PlaceId{},
			},
		},
	}
}

func GetInstance1() instance.Instance {
	return instance.Instance{
		Id:          "conformance_example1",
		PetriNet:    GetPetriNet1().Id,
		TokenCounts: []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}
