package circuit

import (
	"proof-service/workflow"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
)

type TokenChanges struct {
	tokenCountDecreasesPerPlaceId *logderivlookup.Table
	tokenCountDecreasesCount      frontend.Variable
	tokenCountIncreasesPerPlaceId *logderivlookup.Table
	tokenCountIncreasesCount      frontend.Variable
}

type TransitionCircuit struct {
	CurrentInstance   Instance
	CurrentCommitment Commitment
	NextInstance      Instance
	NextCommitment    Commitment
	PetriNet          PetriNet
}

func (circuit *TransitionCircuit) Define(api frontend.API) error {
	api.AssertIsEqual(circuit.NextInstance.PlaceCount, circuit.PetriNet.PlaceCount)
	checkCommitment(api, circuit.CurrentInstance, circuit.CurrentCommitment)
	checkCommitment(api, circuit.NextInstance, circuit.NextCommitment)
	tokenChanges := circuit.computeTokenChanges(api)
	transitionFound := circuit.FindTransition(api, tokenChanges)
	tokenCountsDoNotChange := api.And(api.IsZero(tokenChanges.tokenCountDecreasesCount), api.IsZero(tokenChanges.tokenCountIncreasesCount))
	api.AssertIsEqual(1, api.Or(transitionFound, tokenCountsDoNotChange))
	return nil
}

func (circuit *TransitionCircuit) computeTokenChanges(api frontend.API) TokenChanges {
	tokenCountDecreasesPerPlaceId := logderivlookup.New(api)
	tokenCountIncreasesPerPlaceId := logderivlookup.New(api)
	var tokenCountDecreasesCount frontend.Variable = 0
	var tokenCountIncreasesCount frontend.Variable = 0

	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId]
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId]

		tokenCountStaysTheSame := equals(api, nextTokenCount, currentTokenCount)
		tokenCountDecreases := equals(api, nextTokenCount, api.Sub(currentTokenCount, 1))
		tokenCountIncreases := equals(api, nextTokenCount, api.Add(currentTokenCount, 1))
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

	api.AssertIsLessOrEqual(tokenCountDecreasesCount, workflow.MaxBranchingFactor)
	api.AssertIsLessOrEqual(tokenCountIncreasesCount, workflow.MaxBranchingFactor)

	return TokenChanges{
		tokenCountDecreasesPerPlaceId, tokenCountDecreasesCount, tokenCountIncreasesPerPlaceId, tokenCountIncreasesCount,
	}
}

func (circuit *TransitionCircuit) FindTransition(api frontend.API, tokenChanges TokenChanges) frontend.Variable {

	var transitionFound frontend.Variable = 0

	for i := range circuit.PetriNet.Transitions {
		transition := circuit.PetriNet.Transitions[i]
		matchesIncomingPlaces := equals(api, transition.IncomingPlaceCount, tokenChanges.tokenCountDecreasesCount)
		matchesOutgoingPlaces := equals(api, transition.OutgoingPlaceCount, tokenChanges.tokenCountIncreasesCount)
		// returns 1 for default placeId (workflow.MaxPlaceCount)
		incomingTokenCountsDecrease := tokenChanges.tokenCountDecreasesPerPlaceId.Lookup(transition.IncomingPlaces[:]...)
		outgoingTokenCountsIncrease := tokenChanges.tokenCountIncreasesPerPlaceId.Lookup(transition.OutgoingPlaces[:]...)
		for j := range transition.IncomingPlaces {
			incomingTokenCountDecreases := incomingTokenCountsDecrease[j]
			outgoingTokenCountIncreases := outgoingTokenCountsIncrease[j]
			matchesIncomingPlaces = api.And(matchesIncomingPlaces, incomingTokenCountDecreases)
			matchesOutgoingPlaces = api.And(matchesOutgoingPlaces, outgoingTokenCountIncreases)
		}
		transitionFound = api.Or(transitionFound, api.And(matchesIncomingPlaces, matchesOutgoingPlaces))
	}

	return transitionFound
}
