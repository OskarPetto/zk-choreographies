package testdata

import (
	"proof-service/circuit/input"
	"proof-service/commitment"

	"github.com/consensys/gnark/std/math/uints"
)

func GetCircuitPetriNet1() input.PetriNet {
	return input.PetriNet{
		StartPlace:      uints.NewU8(9),
		PlaceCount:      uints.NewU8(10),
		TransitionCount: uints.NewU8(12),
		Transitions: []input.Transition{
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{9}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{0}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{0}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{8}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{0}),
				OutgoingPlaceCount: uints.NewU8(2),
				OutgoingPlaces:     uints.NewU8Array([]byte{1, 2}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{1}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{5}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{2}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{3}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{3}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{4}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{4}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{2}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(2),
				IncomingPlaces:     uints.NewU8Array([]byte{5, 4}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{6}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(2),
				IncomingPlaces:     uints.NewU8Array([]byte{5, 4}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{7}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{6}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{8}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{7}),
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     uints.NewU8Array([]byte{8}),
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     uints.NewU8Array([]byte{8}),
				OutgoingPlaceCount: uints.NewU8(0),
				OutgoingPlaces:     []uints.U8{},
			},
		},
	}
}

func GetCircuitInstance1() input.Instance {
	return input.Instance{
		PlaceCount:  uints.NewU8(10),
		TokenCounts: uints.NewU8Array([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
	}
}

func GetCircuitCommitment1() input.Commitment {
	return input.Commitment{
		Value:      ([commitment.CommitmentSize]uints.U8)(uints.NewU8Array([]byte{10, 185, 75, 252, 255, 103, 193, 241, 172, 81, 178, 68, 215, 163, 95, 199, 96, 148, 114, 179, 121, 178, 158, 198, 141, 100, 109, 163, 143, 89, 143, 50})),
		Randomness: uints.NewU8Array([]byte{53, 180, 160, 126, 121, 134, 90, 83, 158, 204, 156, 68, 2, 224, 240, 22, 108, 27, 62, 53, 219, 131, 20, 230, 182, 12, 193, 43, 90, 34, 217, 167}),
	}
}
