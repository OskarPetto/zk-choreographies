package proof

import (
	"github.com/consensys/gnark/frontend"
)

const maxPlaceCount = 100
const maxTransitionCount = 100
const maxBranchingFactor = 3

type Instance struct {
	Id                frontend.Variable
	PetriNet          frontend.Variable
	TokenCountsLength frontend.Variable
	TokenCounts       [maxPlaceCount]frontend.Variable
}

type Transition struct {
	IncomingPlaces [maxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaces [maxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	Id          frontend.Variable `gnark:",public"`
	PlaceCount  frontend.Variable `gnark:",public"`
	StartPlace  frontend.Variable `gnark:",public"`
	Transitions [maxTransitionCount]Transition
}
