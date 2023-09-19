package circuit

import (
	"proof-service/workflow"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
)

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
	circuit.checkTokenCounts(api)
	return nil
}

func (circuit *TransitionCircuit) checkTokenCounts(api frontend.API) {
	api.AssertIsEqual(circuit.NextInstance.PlaceCount, circuit.PetriNet.PlaceCount)

	tokenCountDecreasesByPlaceId := logderivlookup.New(api)
	tokenCountIncreasesByPlaceId := logderivlookup.New(api)
	var incomingPlaceCount frontend.Variable = 0
	var outgoingPlaceCount frontend.Variable = 0

	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId]
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId]

		tokenCountStaysTheSame := equals(api, nextTokenCount, currentTokenCount)
		tokenCountDecreases := equals(api, nextTokenCount, api.Sub(currentTokenCount, 1))
		tokenCountIncreases := equals(api, nextTokenCount, api.Add(currentTokenCount, 1))
		api.AssertIsEqual(1, api.Or(api.Or(tokenCountStaysTheSame, tokenCountDecreases), tokenCountIncreases))

		incomingPlaceCount = api.Add(incomingPlaceCount, tokenCountDecreases)
		outgoingPlaceCount = api.Add(outgoingPlaceCount, tokenCountIncreases)

		tokenCountDecreasesByPlaceId.Insert(tokenCountDecreases)
		tokenCountIncreasesByPlaceId.Insert(tokenCountIncreases)

		api.AssertIsBoolean(nextTokenCount)
	}

	// insert 1 at workflow.MaxPlaceCount
	tokenCountDecreasesByPlaceId.Insert(1)
	tokenCountIncreasesByPlaceId.Insert(1)

	api.AssertIsLessOrEqual(incomingPlaceCount, workflow.MaxBranchingFactor)
	api.AssertIsLessOrEqual(outgoingPlaceCount, workflow.MaxBranchingFactor)

	var transitionFound frontend.Variable = 0

	for i := range circuit.PetriNet.Transitions {
		transition := circuit.PetriNet.Transitions[i]
		matchesIncomingPlaces := equals(api, transition.IncomingPlaceCount, incomingPlaceCount)
		matchesOutgoingPlaces := equals(api, transition.OutgoingPlaceCount, outgoingPlaceCount)
		// returns 1 for default placeId (workflow.MaxPlaceCount)
		incomingTokenCountsDecrease := tokenCountDecreasesByPlaceId.Lookup(transition.IncomingPlaces[:]...)
		outgoingTokenCountsIncrease := tokenCountIncreasesByPlaceId.Lookup(transition.OutgoingPlaces[:]...)
		for j := range transition.IncomingPlaces {
			incomingTokenCountDecreases := incomingTokenCountsDecrease[j]
			outgoingTokenCountIncreases := outgoingTokenCountsIncrease[j]
			matchesIncomingPlaces = api.And(matchesIncomingPlaces, incomingTokenCountDecreases)
			matchesOutgoingPlaces = api.And(matchesOutgoingPlaces, outgoingTokenCountIncreases)
		}
		transitionFound = api.Or(transitionFound, api.And(matchesIncomingPlaces, matchesOutgoingPlaces))
	}

	tokenCountsDoNotChange := api.And(api.IsZero(incomingPlaceCount), api.IsZero(outgoingPlaceCount))
	api.AssertIsEqual(1, api.Or(transitionFound, tokenCountsDoNotChange))
}
