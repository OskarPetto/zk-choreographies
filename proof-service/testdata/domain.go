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
		Transitions: []petri_net.Transition{
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{9},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{0},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{0},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{8},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{0},
				OutgoingPlaceCount: 2,
				OutgoingPlaces:     []uint8{1, 2},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{1},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{5},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{2},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{3},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{3},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{4},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{2},
			},
			petri_net.Transition{
				IncomingPlaceCount: 2,
				IncomingPlaces:     []uint8{5, 4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{6},
			},
			petri_net.Transition{
				IncomingPlaceCount: 2,
				IncomingPlaces:     []uint8{5, 4},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{7},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{6},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{8},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{7},
				OutgoingPlaceCount: 1,
				OutgoingPlaces:     []uint8{8},
			},
			petri_net.Transition{
				IncomingPlaceCount: 1,
				IncomingPlaces:     []uint8{8},
				OutgoingPlaceCount: 0,
				OutgoingPlaces:     []uint8{},
			},
		},
	}
}

func GetInstance1() instance.Instance {
	return instance.Instance{
		PlaceCount:  10,
		TokenCounts: []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}

func GetSerializedInstance1() []byte {
	return []byte{10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
}

func GetCommitment1() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example1",
		Value:      [commitment.CommitmentSize]byte{10, 185, 75, 252, 255, 103, 193, 241, 172, 81, 178, 68, 215, 163, 95, 199, 96, 148, 114, 179, 121, 178, 158, 198, 141, 100, 109, 163, 143, 89, 143, 50},
		Randomness: []byte{53, 180, 160, 126, 121, 134, 90, 83, 158, 204, 156, 68, 2, 224, 240, 22, 108, 27, 62, 53, 219, 131, 20, 230, 182, 12, 193, 43, 90, 34, 217, 167},
	}
}
