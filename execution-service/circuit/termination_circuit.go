package circuit

import (
	"github.com/consensys/gnark/frontend"
)

type TerminationCircuit struct {
	ModelHash Hash
	Model     Model
	Instance  Instance
	Signature Signature
}

func (circuit *TerminationCircuit) Define(api frontend.API) error {
	err := checkModelHash(api, circuit.ModelHash, circuit.Model)
	if err != nil {
		return err
	}
	err = checkInstanceHash(api, circuit.Instance)
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
	var atLeastOneEndPlaceHasTokenCountOne frontend.Variable = 0
	for placeId, tokenCount := range circuit.Instance.TokenCounts {
		tokenCountIsOne := equals(api, tokenCount, 1)
		for _, endPlaceId := range circuit.Model.EndPlaces {
			isEndPlace := equals(api, endPlaceId, placeId)
			atLeastOneEndPlaceHasTokenCountOne = api.Or(atLeastOneEndPlaceHasTokenCountOne, api.And(isEndPlace, tokenCountIsOne))
		}
	}
	api.AssertIsEqual(1, atLeastOneEndPlaceHasTokenCountOne)
}
