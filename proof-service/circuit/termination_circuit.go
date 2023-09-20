package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type TerminationCircuit struct {
	Instance   Instance
	Commitment Commitment
	PetriNet   PetriNet
}

func (circuit *TerminationCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.Instance.PlaceCount, circuit.PetriNet.PlaceCount)
	checkCommitment(api, circuit.Instance, circuit.Commitment)
	circuit.checkTokenCounts(api)
	return nil
}

func (circuit *TerminationCircuit) checkTokenCounts(api frontend.API) {
	for placeId := range circuit.Instance.TokenCounts {
		tokenCount := circuit.Instance.TokenCounts[placeId]
		isEndPlace := equals(api, placeId, circuit.PetriNet.EndPlace)
		api.AssertIsEqual(tokenCount, isEndPlace)
	}
}
