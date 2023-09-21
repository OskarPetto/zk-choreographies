package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type InstantiationCircuit struct {
	Instance   Instance
	SaltedHash SaltedHash
	Signature  Signature
	PetriNet   PetriNet
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	err := checkInstanceSaltedHash(api, circuit.Instance, circuit.SaltedHash)
	if err != nil {
		return err
	}
	err = checkSignature(api, circuit.Signature, circuit.SaltedHash)
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
