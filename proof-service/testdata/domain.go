package testdata

import (
	"proof-service/crypto"
	"proof-service/workflow"
)

func GetPetriNet1() workflow.PetriNet {
	return workflow.PetriNet{
		Id:               "conformance_example",
		StartPlace:       9,
		EndPlace:         10,
		PlaceCount:       11,
		ParticipantCount: 1,
		Transitions: []workflow.Transition{
			workflow.Transition{
				Id:                           "As",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{9},
				OutgoingPlaces:               []uint{0},
			},
			workflow.Transition{
				Id:                           "Da1",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{0},
				OutgoingPlaces:               []uint{8},
			},
			workflow.Transition{
				Id:                           "Aa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{0},
				OutgoingPlaces:               []uint{1, 2},
			},
			workflow.Transition{
				Id:                           "Fa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{1},
				OutgoingPlaces:               []uint{5},
			},
			workflow.Transition{
				Id:                           "Sso",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{2},
				OutgoingPlaces:               []uint{3},
			},
			workflow.Transition{
				Id:                           "Ro",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{3},
				OutgoingPlaces:               []uint{4},
			},
			workflow.Transition{
				Id:                           "Co",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{4},
				OutgoingPlaces:               []uint{2},
			},
			workflow.Transition{
				Id:                           "Ao",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{5, 4},
				OutgoingPlaces:               []uint{6},
			},
			workflow.Transition{
				Id:                           "Do",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{5, 4},
				OutgoingPlaces:               []uint{7},
			},
			workflow.Transition{
				Id:                           "Aaa",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{6},
				OutgoingPlaces:               []uint{8},
			},
			workflow.Transition{
				Id:                           "Da2",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{7},
				OutgoingPlaces:               []uint{8},
			},
			workflow.Transition{
				Id:                           "Af",
				IsExecutableByAnyParticipant: true,
				IncomingPlaces:               []uint{8},
				OutgoingPlaces:               []uint{10},
			},
		},
	}
}

func GetPetriNet1Instance1(publicKey []byte) workflow.Instance {
	return workflow.Instance{
		Id:          "conformance_example1",
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetPetriNet1Instance2(publicKey []byte) workflow.Instance {
	return workflow.Instance{
		Id:          "conformance_example2",
		TokenCounts: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetPetriNet1Instance3(publicKey []byte) workflow.Instance {
	return workflow.Instance{
		Id:          "conformance_example4",
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetCommitment1() crypto.SaltedHash {
	return crypto.SaltedHash{
		Value: []byte{15, 119, 210, 82, 4, 149, 235, 173, 255, 201, 90, 205, 146, 233, 251, 58, 54, 88, 10, 179, 75, 101, 147, 46, 127, 239, 221, 252, 28, 71, 138, 66},
		Salt:  []byte{85, 39, 212, 198, 200, 84, 236, 218, 89, 123, 119, 127, 251, 16, 159, 125, 24, 72, 146, 14, 13, 242, 101, 182, 18, 14, 139, 149, 217, 116, 255, 43},
	}
}

func GetPublicKey1() []byte {
	return []byte{70, 200, 160, 220, 129, 215, 38, 174, 106, 10, 190, 160, 109, 87, 219, 147, 161, 184, 34, 209, 190, 54, 152, 202, 123, 230, 254, 52, 193, 43, 56, 147}
}
