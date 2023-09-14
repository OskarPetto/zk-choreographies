package circuit

import (
	"proof-service/circuit/input"
	"proof-service/petri_net"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/uints"
)

type InstantiationCircuit struct {
	Instance   input.Instance
	Commitment input.Commitment
	PetriNet   input.PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.Instance.PlaceCount, circuit.PetriNet.PlaceCount)
	for placeId := 0; placeId < petri_net.MaxPlaceCount; placeId++ {
		isStartPlace := api.IsZero(api.Sub(circuit.PetriNet.StartPlace, placeId))
		expectedTokenCount := api.Select(isStartPlace, 1, 0)
		tokenCount := circuit.Instance.TokenCounts[placeId]
		api.AssertIsEqual(tokenCount, expectedTokenCount)
	}
	hasher, _ := sha2.New(api)
	var hashInput []uints.U8
	hashInput = append(hashInput, circuit.Instance.PlaceCount)
	hashInput = append(hashInput, circuit.Instance.TokenCounts[:]...)
	hashInput = append(hashInput, circuit.Commitment.Randomness[:]...)
	hasher.Write(hashInput)
	hashResult := hasher.Sum()
	api.AssertIsEqual(hashResult, circuit.Commitment.Value)
	return nil
}
