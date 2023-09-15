package testdata

import (
	"proof-service/commitment"
	"proof-service/instance"
	"proof-service/petri_net"
)

func GetPetriNet() petri_net.PetriNet {
	return petri_net.PetriNet{
		StartPlace:      9,
		PlaceCount:      10,
		TransitionCount: 12,
		Transitions: [petri_net.MaxTransitionCount]petri_net.Transition{
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{9},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{0},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{0},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{8},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{0},
				OutgoingPlaceCount: 2,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{1, 2},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{1},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{5},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{2},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{3},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{3},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{4},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{2},
			},
			petri_net.Transition{
				IncomingPlaceCount: 2,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{5, 4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{6},
			},
			petri_net.Transition{
				IncomingPlaceCount: 2,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{5, 4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{7},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{6},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{8},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{7},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{8},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uint8{8},
				OutgoingPlaceCount: 0,
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uint8{},
			},
		},
	}
}

func GetInstance1() instance.Instance {
	return instance.Instance{
		PlaceCount:  10,
		TokenCounts: [petri_net.MaxPlaceCount]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}

func GetSerializedInstance1() []byte {
	bytes := [petri_net.MaxPlaceCount + 1]byte{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	return bytes[:]
}

func GetCommitment1() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example1",
		Value:      [commitment.CommitmentSize]byte{125, 151, 89, 109, 91, 252, 200, 227, 190, 160, 198, 218, 21, 152, 148, 77, 24, 132, 95, 198, 105, 72, 117, 98, 107, 239, 8, 116, 223, 237, 146, 34},
		Randomness: [commitment.RandomnessSize]byte{113, 137, 111, 247, 114, 4, 181, 51, 41, 35, 76, 57, 58, 89, 194, 160, 156, 41, 145, 178, 79, 84, 151, 181, 75, 182, 178, 102, 31, 47, 235, 66},
	}
}
