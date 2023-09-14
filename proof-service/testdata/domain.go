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
	return []byte{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
}

func GetCommitment1() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example1",
		Value:      [commitment.CommitmentSize]byte{173, 250, 76, 199, 33, 125, 245, 112, 193, 231, 245, 111, 241, 47, 83, 227, 59, 97, 235, 27, 41, 252, 154, 174, 203, 125, 46, 134, 192, 60, 21, 61},
		Randomness: []byte{235, 234, 56, 167, 9, 35, 69, 86, 84, 189, 131, 113, 74, 36, 83, 0, 32, 84, 137, 78, 186, 94, 40, 36, 195, 52, 216, 89, 95, 175, 115, 196},
	}
}
