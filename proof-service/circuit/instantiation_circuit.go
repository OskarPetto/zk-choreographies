package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
)

type InstantiationCircuit struct {
	Instance   Instance
	Commitment Commitment
	PetriNet   PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}
	checkCommitment(api, uapi, circuit.Instance, circuit.Commitment)
	uapi.ByteAssertEq(circuit.Instance.PlaceCount, circuit.PetriNet.PlaceCount)
	circuit.checkTokenCounts(api)
	return nil
}

func (circuit *InstantiationCircuit) checkTokenCounts(api frontend.API) {
	for placeId := range circuit.Instance.TokenCounts {
		tokenCount := circuit.Instance.TokenCounts[placeId].Val
		isStartPlace := equals(api, placeId, circuit.PetriNet.StartPlace.Val)
		api.AssertIsEqual(tokenCount, isStartPlace)
	}
}
