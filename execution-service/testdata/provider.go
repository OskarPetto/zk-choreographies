package testdata

import (
	"execution-service/instance"
	"execution-service/model"
)

func GetModel1() model.Model {
	return model.Model{
		Id:         "conformance_example",
		StartPlace: 9,
		PlaceCount: 10,
		Transitions: []model.Transition{
			model.Transition{
				Id:             "As",
				IncomingPlaces: []model.PlaceId{9},
				OutgoingPlaces: []model.PlaceId{0},
			},
			model.Transition{
				Id:             "Da1",
				IncomingPlaces: []model.PlaceId{0},
				OutgoingPlaces: []model.PlaceId{8},
			},
			model.Transition{
				Id:             "Aa",
				IncomingPlaces: []model.PlaceId{0},
				OutgoingPlaces: []model.PlaceId{1, 2},
			},
			model.Transition{
				Id:             "Fa",
				IncomingPlaces: []model.PlaceId{1},
				OutgoingPlaces: []model.PlaceId{5},
			},
			model.Transition{
				Id:             "Sso",
				IncomingPlaces: []model.PlaceId{2},
				OutgoingPlaces: []model.PlaceId{3},
			},
			model.Transition{
				Id:             "Ro",
				IncomingPlaces: []model.PlaceId{3},
				OutgoingPlaces: []model.PlaceId{4},
			},
			model.Transition{
				Id:             "Co",
				IncomingPlaces: []model.PlaceId{4},
				OutgoingPlaces: []model.PlaceId{2},
			},
			model.Transition{
				Id:             "Ao",
				IncomingPlaces: []model.PlaceId{5, 4},
				OutgoingPlaces: []model.PlaceId{6},
			},
			model.Transition{
				Id:             "Do",
				IncomingPlaces: []model.PlaceId{5, 4},
				OutgoingPlaces: []model.PlaceId{7},
			},
			model.Transition{
				Id:             "Aaa",
				IncomingPlaces: []model.PlaceId{6},
				OutgoingPlaces: []model.PlaceId{8},
			},
			model.Transition{
				Id:             "Da2",
				IncomingPlaces: []model.PlaceId{7},
				OutgoingPlaces: []model.PlaceId{8},
			},
			model.Transition{
				Id:             "Af",
				IncomingPlaces: []model.PlaceId{8},
				OutgoingPlaces: []model.PlaceId{},
			},
		},
	}
}

func GetInstance1() instance.Instance {
	return instance.Instance{
		Id:          "conformance_example1",
		Model:       GetModel1().Id,
		TokenCounts: []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}
