package testdata

import (
	"proof-service/model"
)

func GetPetriNet1() model.PetriNet {
	return model.PetriNet{
		Id:               "conformance_example",
		StartPlace:       9,
		EndPlace:         10,
		PlaceCount:       11,
		ParticipantCount: 1,
		Transitions: []model.Transition{
			model.Transition{
				Id:                           "As",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{9},
				OutgoingPlaces:               []uint{0},
			},
			model.Transition{
				Id:                           "Da1",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{0},
				OutgoingPlaces:               []uint{8},
			},
			model.Transition{
				Id:                           "Aa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{0},
				OutgoingPlaces:               []uint{1, 2},
			},
			model.Transition{
				Id:                           "Fa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{1},
				OutgoingPlaces:               []uint{5},
			},
			model.Transition{
				Id:                           "Sso",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{2},
				OutgoingPlaces:               []uint{3},
			},
			model.Transition{
				Id:                           "Ro",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{3},
				OutgoingPlaces:               []uint{4},
			},
			model.Transition{
				Id:                           "Co",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{4},
				OutgoingPlaces:               []uint{2},
			},
			model.Transition{
				Id:                           "Ao",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{5, 4},
				OutgoingPlaces:               []uint{6},
			},
			model.Transition{
				Id:                           "Do",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{5, 4},
				OutgoingPlaces:               []uint{7},
			},
			model.Transition{
				Id:                           "Aaa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{6},
				OutgoingPlaces:               []uint{8},
			},
			model.Transition{
				Id:                           "Da2",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{7},
				OutgoingPlaces:               []uint{8},
			},
			model.Transition{
				Id:                           "Af",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{8},
				OutgoingPlaces:               []uint{10},
			},
		},
	}
}
