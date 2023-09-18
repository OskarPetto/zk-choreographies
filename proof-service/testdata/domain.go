package testdata

import (
	"proof-service/commitment"
	"proof-service/domain"
)

func GetPetriNet1() domain.PetriNet {
	return domain.PetriNet{
		Id:              "conformance_example",
		StartPlace:      9,
		PlaceCount:      10,
		TransitionCount: 12,
		Transitions: [domain.MaxTransitionCount]domain.Transition{
			domain.Transition{
				Id:                 "As",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{9},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{0},
			},
			domain.Transition{
				Id:                 "Da1",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{0},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{8},
			},
			domain.Transition{
				Id:                 "Aa",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{0},
				OutgoingPlaceCount: 2,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{1, 2},
			},
			domain.Transition{
				Id:                 "Fa",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{1},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{5},
			},
			domain.Transition{
				Id:                 "Sso",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{2},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{3},
			},
			domain.Transition{
				Id:                 "Ro",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{3},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{4},
			},
			domain.Transition{
				Id:                 "Co",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{2},
			},
			domain.Transition{
				Id:                 "Ao",
				IncomingPlaceCount: 2,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{5, 4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{6},
			},
			domain.Transition{
				Id:                 "Do",
				IncomingPlaceCount: 2,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{5, 4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{7},
			},
			domain.Transition{
				Id:                 "Aaa",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{6},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{8},
			},
			domain.Transition{
				Id:                 "Da2",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{7},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{8},
			},
			domain.Transition{
				Id:                 "Af",
				IncomingPlaceCount: 1,
				IncomingPlaces:     [domain.MaxBranchingFactor]uint8{8},
				OutgoingPlaceCount: 0,
				OutgoingPlaces:     [domain.MaxBranchingFactor]uint8{},
			},
		},
	}
}

func GetPetriNet1Instance1() domain.Instance {
	return domain.Instance{
		Id:          "conformance_example1",
		PlaceCount:  10,
		TokenCounts: [domain.MaxPlaceCount]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}

func GetPetriNet1Instance2() domain.Instance {
	return domain.Instance{
		Id:          "conformance_example2",
		PlaceCount:  10,
		TokenCounts: [domain.MaxPlaceCount]int8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func GetPetriNet1Instance3() domain.Instance {
	return domain.Instance{
		Id:          "conformance_example3",
		PlaceCount:  10,
		TokenCounts: [domain.MaxPlaceCount]int8{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	}
}

func GetPetriNet1Instance4() domain.Instance {
	return domain.Instance{
		Id:          "conformance_example4",
		PlaceCount:  10,
		TokenCounts: [domain.MaxPlaceCount]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
}

func GetPetriNet1Instance1Serialized() []byte {
	bytes := [domain.MaxPlaceCount + 1]byte{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	return bytes[:]
}

func GetPetriNet1Instance2Serialized() []byte {
	bytes := [domain.MaxPlaceCount + 1]byte{10, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	return bytes[:]
}

func GetPetriNet1Instance3Serialized() []byte {
	bytes := [domain.MaxPlaceCount + 1]byte{10, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	return bytes[:]
}

func GetPetriNet1Instance4Serialized() []byte {
	bytes := [domain.MaxPlaceCount + 1]byte{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	return bytes[:]
}

func GetPetriNet1Instance1Commitment() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example1",
		Value:      [commitment.CommitmentSize]byte{125, 151, 89, 109, 91, 252, 200, 227, 190, 160, 198, 218, 21, 152, 148, 77, 24, 132, 95, 198, 105, 72, 117, 98, 107, 239, 8, 116, 223, 237, 146, 34},
		Randomness: [commitment.RandomnessSize]byte{113, 137, 111, 247, 114, 4, 181, 51, 41, 35, 76, 57, 58, 89, 194, 160, 156, 41, 145, 178, 79, 84, 151, 181, 75, 182, 178, 102, 31, 47, 235, 66},
	}
}

func GetPetriNet1Instance2Commitment() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example2",
		Value:      [commitment.CommitmentSize]byte{58, 122, 231, 60, 219, 204, 89, 72, 126, 138, 51, 142, 50, 143, 204, 186, 133, 217, 136, 138, 243, 90, 185, 72, 92, 6, 82, 43, 203, 128, 159, 159},
		Randomness: [commitment.RandomnessSize]byte{180, 36, 151, 41, 32, 17, 9, 82, 115, 200, 233, 194, 25, 157, 45, 104, 255, 183, 200, 88, 67, 177, 124, 84, 165, 238, 147, 226, 162, 161, 93, 79},
	}
}

func GetPetriNet1Instance3Commitment() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example3",
		Value:      [commitment.CommitmentSize]byte{107, 28, 98, 148, 123, 116, 134, 237, 158, 57, 207, 101, 26, 90, 216, 2, 68, 21, 139, 233, 249, 194, 91, 87, 25, 36, 116, 194, 35, 160, 156, 195},
		Randomness: [commitment.RandomnessSize]byte{8, 153, 6, 247, 21, 48, 228, 49, 164, 188, 199, 154, 13, 209, 38, 77, 81, 178, 246, 77, 87, 55, 172, 228, 90, 198, 184, 57, 227, 41, 90, 205},
	}
}

func GetPetriNet1Instance4Commitment() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example4",
		Value:      [commitment.CommitmentSize]byte{170, 43, 154, 34, 208, 145, 177, 193, 86, 247, 225, 4, 155, 243, 57, 154, 89, 91, 225, 228, 233, 182, 192, 115, 185, 17, 217, 37, 149, 120, 169, 229},
		Randomness: [commitment.RandomnessSize]byte{79, 219, 230, 161, 234, 110, 28, 102, 220, 99, 110, 32, 79, 49, 43, 17, 130, 224, 66, 213, 30, 251, 46, 0, 153, 41, 82, 34, 201, 198, 42, 52},
	}
}
