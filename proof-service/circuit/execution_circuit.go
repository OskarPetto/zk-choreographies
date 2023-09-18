package circuit

import (
	"proof-service/domain"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
	"github.com/consensys/gnark/std/math/uints"
)

type ExecutionCircuit struct {
	CurrentInstance   Instance
	CurrentCommitment Commitment
	NextInstance      Instance
	NextCommitment    Commitment
	PetriNet          PetriNet
}

func (circuit *ExecutionCircuit) Define(api frontend.API) error {
	uapi, err := uints.New[uints.U32](api)
	if err != nil {
		return err
	}
	checkCommitment(api, uapi, circuit.CurrentInstance, circuit.CurrentCommitment)
	checkCommitment(api, uapi, circuit.NextInstance, circuit.NextCommitment)
	api.AssertIsEqual(circuit.NextInstance.PlaceCount.Val, circuit.PetriNet.PlaceCount)
	circuit.checkTokenCounts(api)
	return nil
}

func (circuit *ExecutionCircuit) checkTokenCounts(api frontend.API) {
	tokenCountDecreasesByPlaceId := logderivlookup.New(api)
	tokenCountIncreasesByPlaceId := logderivlookup.New(api)
	var incomingPlaceCount frontend.Variable
	var outgoingPlaceCount frontend.Variable
	incomingPlaceCount = 0
	outgoingPlaceCount = 0
	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId].Val
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId].Val

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

	// insert 1 at domain.MaxPlaceCount
	tokenCountDecreasesByPlaceId.Insert(1)
	tokenCountIncreasesByPlaceId.Insert(1)

	api.AssertIsLessOrEqual(incomingPlaceCount, domain.MaxBranchingFactor)
	api.AssertIsLessOrEqual(outgoingPlaceCount, domain.MaxBranchingFactor)

	var transitionFound frontend.Variable
	transitionFound = 0

	for i := range circuit.PetriNet.Transitions {
		transition := circuit.PetriNet.Transitions[i]
		matchesIncomingPlaces := equals(api, transition.IncomingPlaceCount, incomingPlaceCount)
		matchesOutgoingPlaces := equals(api, transition.OutgoingPlaceCount, outgoingPlaceCount)
		// returns 1 for default placeId (domain.MaxPlaceCount)
		incomingTokenCountsDecrease := tokenCountDecreasesByPlaceId.Lookup(transition.IncomingPlaces[:]...)
		outgoingTokenCountsIncrease := tokenCountIncreasesByPlaceId.Lookup(transition.OutgoingPlaces[:]...)
		for j := 0; j < domain.MaxBranchingFactor; j++ {
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

func greaterThanOrEqual(api frontend.API, a, b frontend.Variable) frontend.Variable {
	cmpResult := api.Cmp(a, b)
	return api.Or(equals(api, cmpResult, 1), api.IsZero(cmpResult))
}
