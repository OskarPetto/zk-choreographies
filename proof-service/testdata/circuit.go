package testdata

import (
	"proof-service/circuit/input"
	"proof-service/commitment"
	"proof-service/petri_net"

	"github.com/consensys/gnark/std/math/uints"
)

func GetCircuitPetriNet1() input.PetriNet {
	return input.PetriNet{
		StartPlace:      uints.NewU8(9),
		PlaceCount:      uints.NewU8(10),
		TransitionCount: uints.NewU8(12),
		Transitions: [petri_net.MaxTransitionCount]input.Transition{
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(9)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(0)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(0)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(0)},
				OutgoingPlaceCount: uints.NewU8(2),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(1), uints.NewU8(2)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(1)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(5)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(2)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(3)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(3)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(4)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(4)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(2)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(2),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(5), uints.NewU8(4)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(6)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(2),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(5), uints.NewU8(4)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(7)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(6)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(7)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
			},
			input.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
				OutgoingPlaceCount: uints.NewU8(0),
				OutgoingPlaces:     [petri_net.MaxBranchingFactor]uints.U8{},
			},
		},
	}
}

func GetCircuitInstance1() input.Instance {
	return input.Instance{
		PlaceCount:  uints.NewU8(10),
		TokenCounts: [petri_net.MaxPlaceCount]uints.U8{uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(1)},
	}
}

func GetCircuitCommitment1() input.Commitment {
	return input.Commitment{
		Value:      [commitment.CommitmentSize]uints.U8{uints.NewU8(125), uints.NewU8(151), uints.NewU8(89), uints.NewU8(109), uints.NewU8(91), uints.NewU8(252), uints.NewU8(200), uints.NewU8(227), uints.NewU8(190), uints.NewU8(160), uints.NewU8(198), uints.NewU8(218), uints.NewU8(21), uints.NewU8(152), uints.NewU8(148), uints.NewU8(77), uints.NewU8(24), uints.NewU8(132), uints.NewU8(95), uints.NewU8(198), uints.NewU8(105), uints.NewU8(72), uints.NewU8(117), uints.NewU8(98), uints.NewU8(107), uints.NewU8(239), uints.NewU8(8), uints.NewU8(116), uints.NewU8(223), uints.NewU8(237), uints.NewU8(146), uints.NewU8(34)},
		Randomness: [commitment.RandomnessSize]uints.U8{uints.NewU8(113), uints.NewU8(137), uints.NewU8(111), uints.NewU8(247), uints.NewU8(114), uints.NewU8(4), uints.NewU8(181), uints.NewU8(51), uints.NewU8(41), uints.NewU8(35), uints.NewU8(76), uints.NewU8(57), uints.NewU8(58), uints.NewU8(89), uints.NewU8(194), uints.NewU8(160), uints.NewU8(156), uints.NewU8(41), uints.NewU8(145), uints.NewU8(178), uints.NewU8(79), uints.NewU8(84), uints.NewU8(151), uints.NewU8(181), uints.NewU8(75), uints.NewU8(182), uints.NewU8(178), uints.NewU8(102), uints.NewU8(31), uints.NewU8(47), uints.NewU8(235), uints.NewU8(66)},
	}
}
