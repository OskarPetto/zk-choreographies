package circuit

import (
	"proof-service/circuit/input"

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
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}
	uapi.ByteAssertEq(circuit.Instance.PlaceCount, circuit.PetriNet.PlaceCount)
	for i := range circuit.Instance.TokenCounts {
		isStartPlace := api.IsZero(api.Sub(circuit.PetriNet.StartPlace.Val, i))
		expectedTokenCount := api.Select(isStartPlace, 1, 0)
		api.AssertIsEqual(circuit.Instance.TokenCounts[i].Val, expectedTokenCount)
	}
	hasher, _ := sha2.New(api)
	hasher.Write([]uints.U8{circuit.Instance.PlaceCount})
	hasher.Write(circuit.Instance.TokenCounts[:])
	hasher.Write(circuit.Commitment.Randomness[:])
	hashResult := hasher.Sum()
	for i := range circuit.Commitment.Value {
		uapi.ByteAssertEq(circuit.Commitment.Value[i], hashResult[i])
	}
	return nil
}
