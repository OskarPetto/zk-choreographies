package testdata

import (
	"proof-service/circuit/input"
	"proof-service/commitment"

	"github.com/consensys/gnark/std/math/uints"
)

func GetCircuitPetriNet1() input.PetriNet {
	return input.PetriNet{
		StartPlace: uints.NewU8(9),
		PlaceCount: uints.NewU8(10),
		Transitions: []input.Transition{
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{9}),
				OutgoingPlaces: uints.NewU8Array([]byte{0}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{0}),
				OutgoingPlaces: uints.NewU8Array([]byte{8}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{0}),
				OutgoingPlaces: uints.NewU8Array([]byte{1, 2}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{1}),
				OutgoingPlaces: uints.NewU8Array([]byte{5}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{2}),
				OutgoingPlaces: uints.NewU8Array([]byte{3}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{3}),
				OutgoingPlaces: uints.NewU8Array([]byte{4}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{4}),
				OutgoingPlaces: uints.NewU8Array([]byte{2}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{5, 4}),
				OutgoingPlaces: uints.NewU8Array([]byte{6}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{5, 4}),
				OutgoingPlaces: uints.NewU8Array([]byte{7}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{6}),
				OutgoingPlaces: uints.NewU8Array([]byte{8}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{7}),
				OutgoingPlaces: uints.NewU8Array([]byte{8}),
			},
			input.Transition{
				IncomingPlaces: uints.NewU8Array([]byte{8}),
				OutgoingPlaces: []uints.U8{},
			},
		},
	}
}

func GetCircuitInstance1() input.Instance {
	return input.Instance{
		TokenCounts: uints.NewU8Array([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}),
	}
}

func GetCircuitCommitment1() input.Commitment {
	return input.Commitment{
		Value:      ([commitment.CommitmentSize]uints.U8)(uints.NewU8Array([]byte{44, 182, 252, 59, 144, 180, 243, 193, 118, 2, 38, 184, 34, 94, 250, 21, 51, 23, 123, 77, 68, 236, 21, 124, 14, 133, 16, 231, 95, 109, 200, 107})),
		Randomness: uints.NewU8Array([]byte{50, 26, 251, 86, 38, 41, 152, 186, 27, 108, 235, 69, 69, 234, 197, 190, 97, 72, 189, 11, 176, 72, 63, 61, 27, 242, 138, 67, 81, 70, 55, 214}),
	}
}
