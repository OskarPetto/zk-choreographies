package circuit

import (
	"execution-service/domain"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash/mimc"
)

type TerminationCircuit struct {
	Model          Model
	Instance       Instance
	Authentication Authentication
	EndPlace       EndPlace
}

func NewTerminationCircuit() TerminationCircuit {
	return TerminationCircuit{
		Authentication: Authentication{
			MerkleProof: merkle.MerkleProof{
				Path: make([]frontend.Variable, domain.MaxParticipantDepth+1),
			},
		},
		EndPlace: EndPlace{
			MerkleProof: merkle.MerkleProof{
				Path: make([]frontend.Variable, domain.MaxEndPlaceDepth+1),
			},
		},
	}
}

func (circuit *TerminationCircuit) Define(api frontend.API) error {
	err := checkModelHash(api, circuit.Model)
	if err != nil {
		return err
	}
	err = checkInstanceHash(api, circuit.Instance)
	if err != nil {
		return err
	}
	checkAuthentication(api, circuit.Authentication, circuit.Instance)
	return circuit.checkTokenCounts(api)
}

func (circuit *TerminationCircuit) checkTokenCounts(api frontend.API) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.EndPlace.MerkleProof.RootHash, circuit.Model.EndPlaceRoot)
	circuit.EndPlace.MerkleProof.VerifyProof(api, &mimc, circuit.EndPlace.Index)
	endPlace := circuit.EndPlace.MerkleProof.Path[0]

	var atLeastOneEndPlaceHasTokenCountOne frontend.Variable = 0
	for placeId, tokenCount := range circuit.Instance.TokenCounts {
		tokenCountIsOne := equals(api, tokenCount, 1)
		isEndPlace := equals(api, placeId, endPlace)
		atLeastOneEndPlaceHasTokenCountOne = api.Or(atLeastOneEndPlaceHasTokenCountOne, api.And(isEndPlace, tokenCountIsOne))
	}
	api.AssertIsEqual(1, atLeastOneEndPlaceHasTokenCountOne)
	return nil
}
