package testdata

import (
	"proof-service/commitment"
	"proof-service/instance"
	"proof-service/petri_net"
)

func GetPetriNet() petri_net.PetriNet {
	return petri_net.PetriNet{
		StartPlace: 9,
		PlaceCount: 10,
		Transitions: []petri_net.Transition{
			petri_net.Transition{
				IncomingPlaces: []uint8{9},
				OutgoingPlaces: []uint8{0},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{0},
				OutgoingPlaces: []uint8{8},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{0},
				OutgoingPlaces: []uint8{1, 2},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{1},
				OutgoingPlaces: []uint8{5},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{2},
				OutgoingPlaces: []uint8{3},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{3},
				OutgoingPlaces: []uint8{4},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{4},
				OutgoingPlaces: []uint8{2},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{5, 4},
				OutgoingPlaces: []uint8{6},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{5, 4},
				OutgoingPlaces: []uint8{7},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{6},
				OutgoingPlaces: []uint8{8},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{7},
				OutgoingPlaces: []uint8{8},
			},
			petri_net.Transition{
				IncomingPlaces: []uint8{8},
				OutgoingPlaces: []uint8{},
			},
		},
	}
}

func GetInstance1() instance.Instance {
	return instance.Instance{
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}

func GetSerializedInstance1() []byte {
	return []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
}

func GetCommitment1() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example1",
		Value:      [commitment.CommitmentSize]byte{44, 182, 252, 59, 144, 180, 243, 193, 118, 2, 38, 184, 34, 94, 250, 21, 51, 23, 123, 77, 68, 236, 21, 124, 14, 133, 16, 231, 95, 109, 200, 107},
		Randomness: []byte{50, 26, 251, 86, 38, 41, 152, 186, 27, 108, 235, 69, 69, 234, 197, 190, 97, 72, 189, 11, 176, 72, 63, 61, 27, 242, 138, 67, 81, 70, 55, 214},
	}
}
