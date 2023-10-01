package circuit

import (
	"execution-service/domain"

	"github.com/consensys/gnark/frontend"
)

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
	findParticipantId(api, circuit.Signature, circuit.Instance)
	circuit.checkTokenCounts(api)
	circuit.checkPublicKeys(api)
	circuit.checkMessageHashes(api)
	return nil
}

func (circuit *InstantiationCircuit) checkTokenCounts(api frontend.API) {
	var outOfBoundsOccurred frontend.Variable = 0
	var numberOfTokenCounts frontend.Variable = 0
	for placeId, tokenCount := range circuit.Instance.TokenCounts {
		var isStartPlace frontend.Variable = 0
		for _, startPlaceId := range circuit.Model.StartPlaces {
			isStartPlace = api.Or(isStartPlace, equals(api, placeId, startPlaceId))
		}
		isTokenCount := equals(api, tokenCount, isStartPlace)
		isOutOfBounds := equals(api, tokenCount, domain.OutOfBoundsTokenCount)
		api.AssertIsDifferent(isOutOfBounds, isTokenCount)

		numberOfTokenCounts = api.Add(numberOfTokenCounts, isTokenCount)

		outOfBoundsOccurred = api.Or(outOfBoundsOccurred, isOutOfBounds)
		api.AssertIsEqual(outOfBoundsOccurred, isOutOfBounds)
	}
	api.AssertIsEqual(numberOfTokenCounts, circuit.Model.PlaceCount)
}

func (circuit *InstantiationCircuit) checkPublicKeys(api frontend.API) {
	var outOfBoundsOccurred frontend.Variable = 0
	var numberOfPublicKeys frontend.Variable = 0
	outOfBoundsPublicKey := fromPublicKey(domain.OutOfBoundsPublicKey())
	for _, publicKey := range circuit.Instance.PublicKeys {
		isOutOfBounds := publicKeyEquals(api, publicKey, outOfBoundsPublicKey)
		isPublicKey := api.IsZero(isOutOfBounds)

		numberOfPublicKeys = api.Add(numberOfPublicKeys, isPublicKey)

		outOfBoundsOccurred = api.Or(outOfBoundsOccurred, isOutOfBounds)
		api.AssertIsEqual(outOfBoundsOccurred, isOutOfBounds)
	}
	api.AssertIsEqual(numberOfPublicKeys, circuit.Model.ParticipantCount)
}

func (circuit *InstantiationCircuit) checkMessageHashes(api frontend.API) {
	var outOfBoundsOccurred frontend.Variable = 0
	var numberOfMessageHashes frontend.Variable = 0
	for _, messageHash := range circuit.Instance.MessageHashes {
		isMessageHash := equals(api, messageHash, emptyMessageHash)
		isOutOfBounds := equals(api, messageHash, outOfBoundsMessageHash)
		api.AssertIsDifferent(isMessageHash, isOutOfBounds)

		numberOfMessageHashes = api.Add(numberOfMessageHashes, isMessageHash)

		outOfBoundsOccurred = api.Or(outOfBoundsOccurred, isOutOfBounds)
		api.AssertIsEqual(outOfBoundsOccurred, isOutOfBounds)
	}
	api.AssertIsEqual(numberOfMessageHashes, circuit.Model.MessageCount)
}
