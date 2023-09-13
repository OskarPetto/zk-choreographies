package proof

import (
	"proof-service/commitment"

	"github.com/consensys/gnark/std/math/uints"
)

type Commitment struct {
	Value      [commitment.CommitmentSize]uints.U8 `gnark:",public"`
	Randomness [commitment.RandomnessSize]uints.U8
}

const MaxPlaceCount = 100

type Instance struct {
	TokenCountsLength uints.U8
	TokenCounts       [MaxPlaceCount]uints.U8
	Commitment        Commitment `gnark:",public"`
}

const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type Transition struct {
	IncomingPlaces [MaxBranchingFactor]uints.U8 `gnark:",public"`
	OutgoingPlaces [MaxBranchingFactor]uints.U8 `gnark:",public"`
}

type PetriNet struct {
	PlaceCount  uints.U8 `gnark:",public"`
	StartPlace  uints.U8 `gnark:",public"`
	Transitions [MaxTransitionCount]Transition
}
