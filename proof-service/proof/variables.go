package proof

import "github.com/consensys/gnark/frontend"

const MaxPlaceCount = 100

type Instance struct {
	TokenCountsLength frontend.Variable
	TokenCounts       [MaxPlaceCount]frontend.Variable
	Randomness        frontend.Variable
	Commitment        frontend.Variable `gnark:",public"`
}

const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type Transition struct {
	IncomingPlaces [MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaces [MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount  frontend.Variable `gnark:",public"`
	StartPlace  frontend.Variable `gnark:",public"`
	Transitions [MaxTransitionCount]Transition
}
