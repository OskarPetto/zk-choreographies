package circuit

import (
	"execution-service/domain"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
)

type InstantiationCircuit struct {
	Model          Model
	Instance       Instance
	Authentication Authentication
}

func NewInstantiationCircuit() InstantiationCircuit {
	return InstantiationCircuit{
		Authentication: Authentication{
			MerkleProof: MerkleProof{
				MerkleProof: merkle.MerkleProof{
					Path: make([]frontend.Variable, domain.MaxParticipantDepth+1),
				},
			},
		},
	}
}

func (circuit *InstantiationCircuit) Define(api frontend.API) error {
	err := checkModelHash(api, circuit.Model, circuit.Instance)
	if err != nil {
		return err
	}
	err = checkInstanceHash(api, circuit.Instance)
	if err != nil {
		return err
	}
	checkAuthentication(api, circuit.Authentication, circuit.Instance)
	circuit.checkTokenCounts(api)
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
