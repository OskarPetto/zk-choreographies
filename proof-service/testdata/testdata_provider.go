package testdata

import (
	"proof-service/commitment"
	"proof-service/proof/inputs"
	"proof-service/workflow"

	"github.com/consensys/gnark/std/math/uints"
)

func GetPetriNet1() workflow.PetriNet {
	return workflow.PetriNet{
		Id:         "conformance_example",
		StartPlace: 9,
		PlaceCount: 10,
		Transitions: []workflow.Transition{
			workflow.Transition{
				Id:             "As",
				IncomingPlaces: []workflow.PlaceId{9},
				OutgoingPlaces: []workflow.PlaceId{0},
			},
			workflow.Transition{
				Id:             "Da1",
				IncomingPlaces: []workflow.PlaceId{0},
				OutgoingPlaces: []workflow.PlaceId{8},
			},
			workflow.Transition{
				Id:             "Aa",
				IncomingPlaces: []workflow.PlaceId{0},
				OutgoingPlaces: []workflow.PlaceId{1, 2},
			},
			workflow.Transition{
				Id:             "Fa",
				IncomingPlaces: []workflow.PlaceId{1},
				OutgoingPlaces: []workflow.PlaceId{5},
			},
			workflow.Transition{
				Id:             "Sso",
				IncomingPlaces: []workflow.PlaceId{2},
				OutgoingPlaces: []workflow.PlaceId{3},
			},
			workflow.Transition{
				Id:             "Ro",
				IncomingPlaces: []workflow.PlaceId{3},
				OutgoingPlaces: []workflow.PlaceId{4},
			},
			workflow.Transition{
				Id:             "Co",
				IncomingPlaces: []workflow.PlaceId{4},
				OutgoingPlaces: []workflow.PlaceId{2},
			},
			workflow.Transition{
				Id:             "Ao",
				IncomingPlaces: []workflow.PlaceId{5, 4},
				OutgoingPlaces: []workflow.PlaceId{6},
			},
			workflow.Transition{
				Id:             "Do",
				IncomingPlaces: []workflow.PlaceId{5, 4},
				OutgoingPlaces: []workflow.PlaceId{7},
			},
			workflow.Transition{
				Id:             "Aaa",
				IncomingPlaces: []workflow.PlaceId{6},
				OutgoingPlaces: []workflow.PlaceId{8},
			},
			workflow.Transition{
				Id:             "Da2",
				IncomingPlaces: []workflow.PlaceId{7},
				OutgoingPlaces: []workflow.PlaceId{8},
			},
			workflow.Transition{
				Id:             "Af",
				IncomingPlaces: []workflow.PlaceId{8},
				OutgoingPlaces: []workflow.PlaceId{},
			},
		},
	}
}

func GetInstance1() workflow.Instance {
	return workflow.Instance{
		Id:          "conformance_example1",
		PetriNet:    GetPetriNet1().Id,
		TokenCounts: []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	}
}

func GetCommitment1() commitment.Commitment {
	return commitment.Commitment{
		Id:         "conformance_example1",
		Value:      [commitment.CommitmentSize]byte{173, 250, 76, 199, 33, 125, 245, 112, 193, 231, 245, 111, 241, 47, 83, 227, 59, 97, 235, 27, 41, 252, 154, 174, 203, 125, 46, 134, 192, 60, 21, 61},
		Randomness: []byte{235, 234, 56, 167, 9, 35, 69, 86, 84, 189, 131, 113, 74, 36, 83, 0, 32, 84, 137, 78, 186, 94, 40, 36, 195, 52, 216, 89, 95, 175, 115, 196},
	}
}

func GetProofPetriNet1() inputs.PetriNet {
	return inputs.PetriNet{
		StartPlace:      uints.NewU8(9),
		PlaceCount:      uints.NewU8(10),
		TransitionCount: uints.NewU8(12),
		Transitions: [inputs.MaxTransitionCount]inputs.Transition{
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(9)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(0)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(0)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(0)},
				OutgoingPlaceCount: uints.NewU8(2),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(1), uints.NewU8(2)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(1)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(5)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(2)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(3)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(3)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(4)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(4)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(2)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(2),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(5), uints.NewU8(4)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(6)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(2),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(5), uints.NewU8(4)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(7)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(6)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(7)},
				OutgoingPlaceCount: uints.NewU8(1),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
			},
			inputs.Transition{
				IncomingPlaceCount: uints.NewU8(1),
				IncomingPlaces:     [inputs.MaxBranchingFactor]uints.U8{uints.NewU8(8)},
				OutgoingPlaceCount: uints.NewU8(0),
				OutgoingPlaces:     [inputs.MaxBranchingFactor]uints.U8{},
			},
		},
	}
}

func GetProofInstance1() inputs.Instance {
	return inputs.Instance{
		PlaceCount:  uints.NewU8(10),
		TokenCounts: [inputs.MaxPlaceCount]uints.U8{uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(0), uints.NewU8(1)},
	}
}

func GetProofCommitment1() inputs.Commitment {
	return inputs.Commitment{
		Value:      [commitment.CommitmentSize]uints.U8{uints.NewU8(173), uints.NewU8(250), uints.NewU8(76), uints.NewU8(199), uints.NewU8(33), uints.NewU8(125), uints.NewU8(245), uints.NewU8(112), uints.NewU8(193), uints.NewU8(231), uints.NewU8(245), uints.NewU8(111), uints.NewU8(241), uints.NewU8(47), uints.NewU8(83), uints.NewU8(227), uints.NewU8(59), uints.NewU8(97), uints.NewU8(235), uints.NewU8(27), uints.NewU8(41), uints.NewU8(252), uints.NewU8(154), uints.NewU8(174), uints.NewU8(203), uints.NewU8(125), uints.NewU8(46), uints.NewU8(134), uints.NewU8(192), uints.NewU8(60), uints.NewU8(21), uints.NewU8(61)},
		Randomness: [...]uints.U8{uints.NewU8(235), uints.NewU8(234), uints.NewU8(56), uints.NewU8(167), uints.NewU8(9), uints.NewU8(35), uints.NewU8(69), uints.NewU8(86), uints.NewU8(84), uints.NewU8(189), uints.NewU8(131), uints.NewU8(113), uints.NewU8(74), uints.NewU8(36), uints.NewU8(83), uints.NewU8(0), uints.NewU8(32), uints.NewU8(84), uints.NewU8(137), uints.NewU8(78), uints.NewU8(186), uints.NewU8(94), uints.NewU8(40), uints.NewU8(36), uints.NewU8(195), uints.NewU8(52), uints.NewU8(216), uints.NewU8(89), uints.NewU8(95), uints.NewU8(175), uints.NewU8(115), uints.NewU8(196)},
	}
}
