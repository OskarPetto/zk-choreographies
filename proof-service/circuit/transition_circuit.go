package circuit

import (
	"proof-service/workflow"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
)

type TokenChanges struct {
	tokenCountDecreasesPerPlaceId *logderivlookup.Table
	tokenCountDecreaseCount       frontend.Variable
	tokenCountIncreasesPerPlaceId *logderivlookup.Table
	tokenCountIncreaseCount       frontend.Variable
}

type TransitionCircuit struct {
	CurrentInstance   Instance
	CurrentCommitment Commitment
	NextInstance      Instance
	NextCommitment    Commitment
	PetriNet          PetriNet
}

func (circuit *TransitionCircuit) Define(api frontend.API) error {
	checkCommitment(api, circuit.CurrentInstance, circuit.CurrentCommitment)
	checkCommitment(api, circuit.NextInstance, circuit.NextCommitment)
	tokenChanges := circuit.checkTokenCounts(api)
	transitionFound := circuit.FindTransition(api, tokenChanges)
	tokenCountsDoNotChange := api.And(api.IsZero(tokenChanges.tokenCountDecreaseCount), api.IsZero(tokenChanges.tokenCountIncreaseCount))
	api.AssertIsEqual(1, api.Or(transitionFound, tokenCountsDoNotChange))
	return nil
}

func (circuit *TransitionCircuit) checkTokenCounts(api frontend.API) TokenChanges {
	api.AssertIsEqual(circuit.NextInstance.PlaceCount, circuit.PetriNet.PlaceCount)

	tokenCountDecreasesPerPlaceId := logderivlookup.New(api)
	tokenCountIncreasesPerPlaceId := logderivlookup.New(api)
	var tokenCountDecreaseCount frontend.Variable = 0
	var tokenCountIncreaseCount frontend.Variable = 0

	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId]
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId]

		tokenCountStaysTheSame := equals(api, nextTokenCount, currentTokenCount)
		tokenCountDecreases := equals(api, nextTokenCount, api.Sub(currentTokenCount, 1))
		tokenCountIncreases := equals(api, nextTokenCount, api.Add(currentTokenCount, 1))
		api.AssertIsEqual(1, api.Or(api.Or(tokenCountStaysTheSame, tokenCountDecreases), tokenCountIncreases))

		tokenCountDecreaseCount = api.Add(tokenCountDecreaseCount, tokenCountDecreases)
		tokenCountIncreaseCount = api.Add(tokenCountIncreaseCount, tokenCountIncreases)

		tokenCountDecreasesPerPlaceId.Insert(tokenCountDecreases)
		tokenCountIncreasesPerPlaceId.Insert(tokenCountIncreases)

		api.AssertIsBoolean(nextTokenCount)
	}

	// insert 1 at workflow.MaxPlaceCount
	tokenCountDecreasesPerPlaceId.Insert(1)
	tokenCountIncreasesPerPlaceId.Insert(1)

	api.AssertIsLessOrEqual(tokenCountDecreaseCount, workflow.MaxBranchingFactor)
	api.AssertIsLessOrEqual(tokenCountIncreaseCount, workflow.MaxBranchingFactor)

	return TokenChanges{
		tokenCountDecreasesPerPlaceId, tokenCountDecreaseCount, tokenCountIncreasesPerPlaceId, tokenCountIncreaseCount,
	}
}

func (circuit *TransitionCircuit) FindTransition(api frontend.API, tokenChanges TokenChanges) frontend.Variable {

	var transitionFound frontend.Variable = 0

	for i := range circuit.PetriNet.Transitions {
		transition := circuit.PetriNet.Transitions[i]
		matchesIncomingPlaces := equals(api, transition.IncomingPlaceCount, tokenChanges.tokenCountDecreaseCount)
		matchesOutgoingPlaces := equals(api, transition.OutgoingPlaceCount, tokenChanges.tokenCountIncreaseCount)
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
