package proof

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type InstantiationCircuit struct {
	Instance Instance
	PetriNet PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.Instance.TokenCountsLength, circuit.PetriNet.PlaceCount)
	for placeId := 0; placeId < MaxPlaceCount; placeId++ {
		isStartPlace := api.IsZero(api.Sub(circuit.PetriNet.StartPlace, placeId))
		expectedTokenCount := api.Select(isStartPlace, 1, 0)
		api.AssertIsEqual(circuit.Instance.TokenCounts[placeId], expectedTokenCount)
	}
	hasher, _ := mimc.NewMiMC(api)
	hasher.Write(circuit.Instance.TokenCountsLength, circuit.Instance.TokenCounts, circuit.Instance.Randomness)
	return nil
}
