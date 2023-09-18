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
		TokenCounts: [domain.MaxPlaceCount]int8{0, 0, 0, 0, 1, 1, 0, 0, 0, 0},
	}
}

func GetPetriNet1Instance4() domain.Instance {
	return domain.Instance{
		Id:          "conformance_example4",
		PlaceCount:  10,
		TokenCounts: [domain.MaxPlaceCount]int8{0, 0, 0, 0, 0, 0, 0, 1, 0, 0},
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
	bytes := [domain.MaxPlaceCount + 1]byte{10, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0}
	return bytes[:]
}

func GetPetriNet1Instance4Serialized() []byte {
	bytes := [domain.MaxPlaceCount + 1]byte{10, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0}
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
		Value:      [commitment.CommitmentSize]byte{137, 32, 200, 146, 163, 67, 93, 114, 36, 1, 139, 191, 112, 105, 88, 224, 253, 193, 207, 59, 167, 184, 245, 199, 106, 126, 162, 83, 163, 246, 116, 242},
		Randomness: [commitment.RandomnessSize]byte{131, 143, 136, 102, 120, 73, 122, 153, 125, 208, 6, 6, 208, 23, 4, 26, 91, 42, 63, 137, 106, 212, 15, 43, 252, 114, 168, 6, 2, 222, 218, 244},
	}
}

func GetPetriNet1Instance4Commitment() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example4",
		Value:      [commitment.CommitmentSize]byte{227, 215, 143, 232, 5, 168, 116, 206, 196, 140, 135, 47, 231, 194, 99, 34, 234, 162, 111, 233, 53, 96, 167, 252, 128, 218, 89, 242, 20, 86, 74, 116},
		Randomness: [commitment.RandomnessSize]byte{83, 247, 198, 223, 212, 151, 241, 171, 214, 252, 166, 173, 39, 96, 9, 69, 51, 42, 138, 88, 230, 202, 148, 54, 95, 239, 77, 29, 134, 108, 195, 135},
	}
}
