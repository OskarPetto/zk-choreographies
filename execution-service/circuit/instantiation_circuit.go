package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
)

type InstantiationCircuit struct {
	Instance  Instance
	Signature Signature
	Model     Model
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
	return circuit.checkMessageHashes(api)
}

func (circuit *InstantiationCircuit) checkTokenCounts(api frontend.API) {
	for placeId, tokenCount := range circuit.Instance.TokenCounts {
		var isStartPlace frontend.Variable = 0
		for _, startPlaceId := range circuit.Model.StartPlaces {
			isStartPlace = api.Or(isStartPlace, equals(api, placeId, startPlaceId))
		}
		api.AssertIsEqual(tokenCount, isStartPlace)
	}
}

func (circuit *InstantiationCircuit) checkMessageHashes(api frontend.API) error {
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}
	for _, messageHash := range circuit.Instance.MessageHashes {
		for _, messageHashByte := range messageHash.Value {
			uapi.ByteAssertEq(messageHashByte, uints.NewU8(0))
		}
	}
	return nil
}
