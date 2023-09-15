package testdata

import "proof-service/workflow"

func GetWorkflowPetriNet1() workflow.PetriNet {
	return workflow.PetriNet{
		Id:         "conformance_example",
		StartPlace: 9,
		PlaceCount: 10,
		Transitions: []workflow.Transition{
			workflow.Transition{
				Id:             "As",
				IncomingPlaces: []workflow.PlaceId{9},
				OutgoingPlaces: []workflow.PlaceId{0},
			},
			workflow.Transition{
				Id:             "Da1",
				IncomingPlaces: []workflow.PlaceId{0},
				OutgoingPlaces: []workflow.PlaceId{8},
			},
			workflow.Transition{
				Id:             "Aa",
				IncomingPlaces: []workflow.PlaceId{0},
				OutgoingPlaces: []workflow.PlaceId{1, 2},
			},
			workflow.Transition{
				Id:             "Fa",
				IncomingPlaces: []workflow.PlaceId{1},
				OutgoingPlaces: []workflow.PlaceId{5},
			},
			workflow.Transition{
				Id:             "Sso",
				IncomingPlaces: []workflow.PlaceId{2},
				OutgoingPlaces: []workflow.PlaceId{3},
			},
			workflow.Transition{
				Id:             "Ro",
				IncomingPlaces: []workflow.PlaceId{3},
				OutgoingPlaces: []workflow.PlaceId{4},
			},
			workflow.Transition{
				Id:             "Co",
				IncomingPlaces: []workflow.PlaceId{4},
				OutgoingPlaces: []workflow.PlaceId{2},
			},
			workflow.Transition{
				Id:             "Ao",
				IncomingPlaces: []workflow.PlaceId{5, 4},
				OutgoingPlaces: []workflow.PlaceId{6},
			},
			workflow.Transition{
				Id:             "Do",
				IncomingPlaces: []workflow.PlaceId{5, 4},
				OutgoingPlaces: []workflow.PlaceId{7},
			},
			workflow.Transition{
				Id:             "Aaa",
				IncomingPlaces: []workflow.PlaceId{6},
				OutgoingPlaces: []workflow.PlaceId{8},
			},
			workflow.Transition{
				Id:             "Da2",
				IncomingPlaces: []workflow.PlaceId{7},
				OutgoingPlaces: []workflow.PlaceId{8},
			},
			workflow.Transition{
				Id:             "Af",
				IncomingPlaces: []workflow.PlaceId{8},
				OutgoingPlaces: []workflow.PlaceId{},
			},
		},
	}
}

func GetWorkflowInstance1() workflow.Instance {
	return workflow.Instance{
		Id:          "conformance_example1",
		PetriNet:    GetWorkflowPetriNet1().Id,
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}
