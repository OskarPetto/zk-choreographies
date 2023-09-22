package circuit

import (
	"proof-service/model"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
)

type TokenCountChanges struct {
	tokenCountDecreasesPerPlaceId *logderivlookup.Table
	tokenCountDecreasesCount      frontend.Variable
	tokenCountIncreasesPerPlaceId *logderivlookup.Table
	tokenCountIncreasesCount      frontend.Variable
	noChanges                     frontend.Variable
}

type TransitionCircuit struct {
	CurrentInstance       Instance
	NextInstance          Instance
	NextInstanceSignature Signature
	PetriNet              PetriNet
}

func (circuit *TransitionCircuit) Define(api frontend.API) error {
	err := checkInstanceHash(api, circuit.CurrentInstance)
	if err != nil {
		return err
	}
	err = checkInstanceHash(api, circuit.NextInstance)
	if err != nil {
		return err
	}
	err = checkSignature(api, circuit.NextInstanceSignature, circuit.NextInstance)
	if err != nil {
		return err
	}
	circuit.comparePublicKeys(api)

	participantId := findParticipantId(api, circuit.NextInstanceSignature, circuit.NextInstance)

	tokenCountChanges := circuit.compareTokenCounts(api)
	circuit.checkTransition(api, tokenCountChanges, participantId)
	return nil
}

func (circuit *TransitionCircuit) comparePublicKeys(api frontend.API) {
	for i := range circuit.CurrentInstance.PublicKeys {
		currentPublicKey := circuit.CurrentInstance.PublicKeys[i]
		nextPublicKey := circuit.NextInstance.PublicKeys[i]
		api.AssertIsEqual(1, publicKeyEquals(api, currentPublicKey, nextPublicKey))
	}
}

func (circuit *TransitionCircuit) compareTokenCounts(api frontend.API) TokenCountChanges {
	tokenCountDecreasesPerPlaceId := logderivlookup.New(api)
	tokenCountIncreasesPerPlaceId := logderivlookup.New(api)
	var tokenCountDecreasesCount frontend.Variable = 0
	var tokenCountIncreasesCount frontend.Variable = 0

	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId]
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId]

		tokenChange := api.Sub(nextTokenCount, currentTokenCount)
		tokenCountStaysTheSame := api.IsZero(tokenChange)
		tokenCountDecreases := equals(api, tokenChange, -1)
		tokenCountIncreases := equals(api, tokenChange, 1)
		api.AssertIsEqual(1, api.Or(api.Or(tokenCountStaysTheSame, tokenCountDecreases), tokenCountIncreases))

		tokenCountDecreasesCount = api.Add(tokenCountDecreasesCount, tokenCountDecreases)
		tokenCountIncreasesCount = api.Add(tokenCountIncreasesCount, tokenCountIncreases)

		tokenCountDecreasesPerPlaceId.Insert(tokenCountDecreases)
		tokenCountIncreasesPerPlaceId.Insert(tokenCountIncreases)

		api.AssertIsBoolean(nextTokenCount)
	}

	// insert 1 at workflow.MaxPlaceCount (default value of incomingPlaces and outgoingPlaces arrays)
	tokenCountDecreasesPerPlaceId.Insert(1)
	tokenCountIncreasesPerPlaceId.Insert(1)

	api.AssertIsLessOrEqual(tokenCountDecreasesCount, model.MaxBranchingFactor)
	api.AssertIsLessOrEqual(tokenCountIncreasesCount, model.MaxBranchingFactor)

	noChanges := api.And(api.IsZero(tokenCountDecreasesCount), api.IsZero(tokenCountIncreasesCount))

	return TokenCountChanges{
		tokenCountDecreasesPerPlaceId, tokenCountDecreasesCount, tokenCountIncreasesPerPlaceId, tokenCountIncreasesCount, noChanges,
	}
}

func (circuit *TransitionCircuit) checkTransition(api frontend.API, tokenCountChanges TokenCountChanges, participantId frontend.Variable) {

	var transitionFound frontend.Variable = 0

	for _, transition := range circuit.PetriNet.Transitions {
		participantIdMatches := api.Or(transition.IsExecutableByAnyParticipant, equals(api, transition.Participant, participantId))
		incomingPlacesMatch := equals(api, transition.IncomingPlaceCount, tokenCountChanges.tokenCountDecreasesCount)
		outgoingPlacesMatch := equals(api, transition.OutgoingPlaceCount, tokenCountChanges.tokenCountIncreasesCount)
		// returns 1 for default placeId (workflow.MaxPlaceCount)
		incomingTokenCountsDecrease := tokenCountChanges.tokenCountDecreasesPerPlaceId.Lookup(transition.IncomingPlaces[:]...)
		outgoingTokenCountsIncrease := tokenCountChanges.tokenCountIncreasesPerPlaceId.Lookup(transition.OutgoingPlaces[:]...)
		for j := range transition.IncomingPlaces {
			incomingTokenCountDecreases := incomingTokenCountsDecrease[j]
			outgoingTokenCountIncreases := outgoingTokenCountsIncrease[j]
			incomingPlacesMatch = api.And(incomingPlacesMatch, incomingTokenCountDecreases)
			outgoingPlacesMatch = api.And(outgoingPlacesMatch, outgoingTokenCountIncreases)
		}
		transitionMatches := api.And(participantIdMatches, api.And(incomingPlacesMatch, outgoingPlacesMatch))
		transitionFound = api.Or(transitionFound, transitionMatches)
	}

	api.AssertIsEqual(1, api.Or(transitionFound, tokenCountChanges.noChanges))
}
