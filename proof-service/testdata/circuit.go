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
		Value:      [commitment.CommitmentSize]uints.U8{uints.NewU8(173), uints.NewU8(250), uints.NewU8(76), uints.NewU8(199), uints.NewU8(33), uints.NewU8(125), uints.NewU8(245), uints.NewU8(112), uints.NewU8(193), uints.NewU8(231), uints.NewU8(245), uints.NewU8(111), uints.NewU8(241), uints.NewU8(47), uints.NewU8(83), uints.NewU8(227), uints.NewU8(59), uints.NewU8(97), uints.NewU8(235), uints.NewU8(27), uints.NewU8(41), uints.NewU8(252), uints.NewU8(154), uints.NewU8(174), uints.NewU8(203), uints.NewU8(125), uints.NewU8(46), uints.NewU8(134), uints.NewU8(192), uints.NewU8(60), uints.NewU8(21), uints.NewU8(61)},
		Randomness: [...]uints.U8{uints.NewU8(235), uints.NewU8(234), uints.NewU8(56), uints.NewU8(167), uints.NewU8(9), uints.NewU8(35), uints.NewU8(69), uints.NewU8(86), uints.NewU8(84), uints.NewU8(189), uints.NewU8(131), uints.NewU8(113), uints.NewU8(74), uints.NewU8(36), uints.NewU8(83), uints.NewU8(0), uints.NewU8(32), uints.NewU8(84), uints.NewU8(137), uints.NewU8(78), uints.NewU8(186), uints.NewU8(94), uints.NewU8(40), uints.NewU8(36), uints.NewU8(195), uints.NewU8(52), uints.NewU8(216), uints.NewU8(89), uints.NewU8(95), uints.NewU8(175), uints.NewU8(115), uints.NewU8(196)},
	}
}
