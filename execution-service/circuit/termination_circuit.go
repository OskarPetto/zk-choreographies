package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type TerminationCircuit struct {
	Instance   Instance
	SaltedHash SaltedHash
	Signature  Signature
	PetriNet   PetriNet
}

func (circuit *TerminationCircuit) Define(api frontend.API) error {
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

func (circuit *TerminationCircuit) checkTokenCounts(api frontend.API) {
	for placeId := range circuit.Instance.TokenCounts {
		tokenCount := circuit.Instance.TokenCounts[placeId]
		isEndPlace := equals(api, placeId, circuit.PetriNet.EndPlace)
		api.AssertIsEqual(tokenCount, isEndPlace)
	}
}
