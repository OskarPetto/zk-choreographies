package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type InstantiationCircuit struct {
	Instance   Instance
	Commitment Commitment
	PetriNet   PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	checkCommitment(api, circuit.Instance, circuit.Commitment)
	circuit.checkTokenCounts(api)
	return nil
}

func (circuit *InstantiationCircuit) checkTokenCounts(api frontend.API) {
	api.AssertIsEqual(circuit.Instance.PlaceCount, circuit.PetriNet.PlaceCount)
	for placeId := range circuit.Instance.TokenCounts {
		tokenCount := circuit.Instance.TokenCounts[placeId]
		isStartPlace := equals(api, placeId, circuit.PetriNet.StartPlace)
		api.AssertIsEqual(tokenCount, isStartPlace)
	}
}
