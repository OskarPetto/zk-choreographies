package proof

import (
	"github.com/consensys/gnark/frontend"
)

type InstantiationCircuit struct {
	Instance Instance
	PetriNet PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.Instance.PetriNet, circuit.PetriNet.Id)
	api.AssertIsEqual(circuit.Instance.TokenCountsLength, circuit.PetriNet.PlaceCount)
	for placeId := 0; placeId < maxPlaceCount; placeId++ {
		isStartPlace := api.IsZero(api.Sub(circuit.PetriNet.StartPlace, placeId))
		expectedTokenCount := api.Select(isStartPlace, 1, 0)
		api.AssertIsEqual(circuit.Instance.TokenCounts[placeId], expectedTokenCount)
	}

	return nil
}
