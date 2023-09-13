package proof

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/uints"
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
		tokenCount := circuit.Instance.TokenCounts[placeId]
		api.AssertIsEqual(tokenCount, expectedTokenCount)
	}
	hasher, _ := sha2.New(api)
	var hashInput []uints.U8
	hashInput = append(hashInput, circuit.Instance.TokenCountsLength)
	hashInput = append(hashInput, circuit.Instance.TokenCounts[:]...)
	hashInput = append(hashInput, circuit.Instance.Commitment.Randomness[:]...)
	hasher.Write(hashInput)
	hashResult := hasher.Sum()
	api.AssertIsEqual(hashResult, circuit.Instance.Commitment.Value)
	return nil
}
