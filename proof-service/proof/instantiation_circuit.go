package proof

import (
	"github.com/consensys/gnark/frontend"
)

const maxPlaceCount = 100
const maxTransitionCount = 100
const maxBranchingFactor = 3

type Instance struct {
	TokenCountsLength frontend.Variable `gnark:",public"`
	TokenCounts       [maxPlaceCount]frontend.Variable
}

type Transition struct {
	IncomingPlaces [maxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaces [maxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount  frontend.Variable `gnark:",public"`
	StartPlace  frontend.Variable `gnark:",public"`
	Transitions [maxTransitionCount]Transition
}

type InstantiationCircuit struct {
	Instance Instance
	PetriNet PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.Instance.TokenCountsLength, circuit.PetriNet.PlaceCount)
	for i := 0; i < maxPlaceCount; i++ {
		tokenCount := circuit.Instance.TokenCounts[i]
		api.AssertIsBoolean(tokenCount)
	}

	return nil
}

func isZeroOrOne(api frontend.API, v frontend.Variable) frontend.Variable {
	return api.Or(api.IsZero(v), api.IsZero(api.Sub(1, v)))
}
