package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type InstantiationCircuit struct {
	Instance  Instance
	Signature Signature
	PetriNet  PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	err := checkInstanceHash(api, circuit.Instance)
	if err != nil {
		return err
	}
	err = checkSignature(api, circuit.Signature, circuit.Instance)
	if err != nil {
		return err
	}
	findParticipantId(api, circuit.Signature, circuit.Instance)
	circuit.checkTokenCounts(api)
	return nil
}

func (circuit *InstantiationCircuit) checkTokenCounts(api frontend.API) {
	for placeId := range circuit.Instance.TokenCounts {
		tokenCount := circuit.Instance.TokenCounts[placeId]
		isStartPlace := equals(api, placeId, circuit.PetriNet.StartPlace)
		api.AssertIsEqual(tokenCount, isStartPlace)
	}
}
