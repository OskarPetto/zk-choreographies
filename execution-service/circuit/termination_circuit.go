package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type TerminationCircuit struct {
	Instance  Instance
	Signature Signature
	PetriNet  PetriNet
}

func (circuit *TerminationCircuit) Define(api frontend.API) error {
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

func (circuit *TerminationCircuit) checkTokenCounts(api frontend.API) {
	for placeId := range circuit.Instance.TokenCounts {
		tokenCount := circuit.Instance.TokenCounts[placeId]
		isEndPlace := equals(api, placeId, circuit.PetriNet.EndPlace)
		api.AssertIsEqual(tokenCount, isEndPlace)
	}
}
