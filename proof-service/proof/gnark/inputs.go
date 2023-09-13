package gnark

import (
	"proof-service/proof"

	"github.com/consensys/gnark/frontend"
)

type Instance struct {
	TokenCountsLength frontend.Variable
	TokenCounts       [proof.MaxPlaceCount]frontend.Variable
	Randomness        frontend.Variable
	Commitment        frontend.Variable `gnark:",public"`
}

type Transition struct {
	IncomingPlaces [proof.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaces [proof.MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount  frontend.Variable `gnark:",public"`
	StartPlace  frontend.Variable `gnark:",public"`
	Transitions [proof.MaxTransitionCount]Transition
}
