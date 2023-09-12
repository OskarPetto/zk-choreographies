package testdata

import "execution-service/model"

type TestdataProvider interface {
	GetModel1() model.Model
}

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
				Id:             "As",
				IncomingPlaces: []model.PlaceId{9},
				OutgoingPlaces: []model.PlaceId{0},
			},
			model.Transition{
				Id:             "As",
				IncomingPlaces: []model.PlaceId{9},
				OutgoingPlaces: []model.PlaceId{0},
			},
			model.Transition{
				Id:             "As",
				IncomingPlaces: []model.PlaceId{9},
				OutgoingPlaces: []model.PlaceId{0},
			},
			model.Transition{
				Id:             "As",
				IncomingPlaces: []model.PlaceId{9},
				OutgoingPlaces: []model.PlaceId{0},
			},
			model.Transition{
				Id:             "As",
				IncomingPlaces: []model.PlaceId{9},
				OutgoingPlaces: []model.PlaceId{0},
			},
			model.Transition{
				Id:             "As",
				IncomingPlaces: []model.PlaceId{9},
				OutgoingPlaces: []model.PlaceId{0},
			},
		},
	}
}
