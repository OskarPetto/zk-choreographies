package testdata

import (
	"proof-service/domain"
)

func GetPetriNet1() domain.PetriNet {
	return domain.PetriNet{
		Id:               "conformance_example",
		StartPlace:       9,
		EndPlace:         10,
		PlaceCount:       11,
		ParticipantCount: 1,
		Transitions: []domain.Transition{
			domain.Transition{
				Id:                           "As",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{9},
				OutgoingPlaces:               []uint{0},
			},
			domain.Transition{
				Id:                           "Da1",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{0},
				OutgoingPlaces:               []uint{8},
			},
			domain.Transition{
				Id:                           "Aa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{0},
				OutgoingPlaces:               []uint{1, 2},
			},
			domain.Transition{
				Id:                           "Fa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{1},
				OutgoingPlaces:               []uint{5},
			},
			domain.Transition{
				Id:                           "Sso",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{2},
				OutgoingPlaces:               []uint{3},
			},
			domain.Transition{
				Id:                           "Ro",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{3},
				OutgoingPlaces:               []uint{4},
			},
			domain.Transition{
				Id:                           "Co",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{4},
				OutgoingPlaces:               []uint{2},
			},
			domain.Transition{
				Id:                           "Ao",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{5, 4},
				OutgoingPlaces:               []uint{6},
			},
			domain.Transition{
				Id:                           "Do",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{5, 4},
				OutgoingPlaces:               []uint{7},
			},
			domain.Transition{
				Id:                           "Aaa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{6},
				OutgoingPlaces:               []uint{8},
			},
			domain.Transition{
				Id:                           "Da2",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{7},
				OutgoingPlaces:               []uint{8},
			},
			domain.Transition{
				Id:                           "Af",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{8},
				OutgoingPlaces:               []uint{10},
			},
		},
	}
}
