package circuit

import (
	"execution-service/domain"

	"github.com/consensys/gnark/frontend"
)

const defaultMessageHash = 0
const invalidMessageHash = 1

type InstantiationCircuit struct {
	Model     Model
	Instance  Instance
	Signature Signature
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	err := checkModelHash(api, circuit.Model)
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
	api.AssertIsLessOrEqual(circuit.Model.CreatedAt, circuit.Instance.CreatedAt)
	findParticipantId(api, circuit.Signature, circuit.Instance)
	circuit.checkTokenCounts(api)
	circuit.checkMessageHashes(api)
	return nil
}

func (circuit *InstantiationCircuit) checkTokenCounts(api frontend.API) {
	var invalidTokenCountOccurred frontend.Variable = 0
	var numberOfValidTokenCounts frontend.Variable = 0
	for placeId, tokenCount := range circuit.Instance.TokenCounts {
		var isStartPlace frontend.Variable = 0
		for _, startPlaceId := range circuit.Model.StartPlaces {
			isStartPlace = api.Or(isStartPlace, equals(api, placeId, startPlaceId))
		}
		tokenCountIsValid := equals(api, tokenCount, isStartPlace)
		tokenCountIsInValid := equals(api, tokenCount, domain.InvalidTokenCount)
		api.AssertIsDifferent(tokenCountIsInValid, tokenCountIsValid)

		numberOfValidTokenCounts = api.Add(numberOfValidTokenCounts, tokenCountIsValid)

		invalidTokenCountOccurred = api.Or(invalidTokenCountOccurred, tokenCountIsInValid)
		api.AssertIsEqual(invalidTokenCountOccurred, tokenCountIsInValid)
	}
	api.AssertIsEqual(numberOfValidTokenCounts, circuit.Model.PlaceCount)
}

func (circuit *InstantiationCircuit) checkMessageHashes(api frontend.API) {
	var invalidMessageHashOccurred frontend.Variable = 0
	var numberOfValidMessageHashes frontend.Variable = 0
	for _, messageHash := range circuit.Instance.MessageHashes {
		messageHashIsValid := equals(api, messageHash, defaultMessageHash)
		messageHashIsInvalid := equals(api, messageHash, invalidMessageHash)
		api.AssertIsDifferent(messageHashIsValid, messageHashIsInvalid)

		numberOfValidMessageHashes = api.Add(numberOfValidMessageHashes, messageHashIsValid)

		invalidMessageHashOccurred = api.Or(invalidMessageHashOccurred, messageHashIsInvalid)
		api.AssertIsEqual(invalidMessageHashOccurred, messageHashIsInvalid)
	}
	api.AssertIsEqual(numberOfValidMessageHashes, circuit.Model.MessageCount)
}
