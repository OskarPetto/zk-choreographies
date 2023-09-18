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

	api.AssertIsLessOrEqual(incomingPlaceCount, api.Sub(domain.MaxBranchingFactor, 1))
	api.AssertIsLessOrEqual(outgoingPlaceCount, api.Sub(domain.MaxBranchingFactor, 1))

	var transitionFound frontend.Variable
	transitionFound = 0

	for i := range circuit.PetriNet.Transitions {
		transition := circuit.PetriNet.Transitions[i]
		var matchesIncomingPlaces frontend.Variable
		matchesIncomingPlaces = equals(api, transition.IncomingPlaceCount, incomingPlaceCount)
		tokenCountDecreasesOfTransition := tokenCountDecreasesByPlaceId.Lookup(transition.IncomingPlaces[:]...)
		for j := range tokenCountDecreasesOfTransition {
			isPlaceIdInvalid := greaterThanOrEqual(api, j, transition.IncomingPlaceCount)
			tokenCountDecreases := tokenCountDecreasesOfTransition[j]
			matchesIncomingPlaces = api.And(matchesIncomingPlaces, api.Or(isPlaceIdInvalid, tokenCountDecreases))
		}
		var matchesOutgoingPlaces frontend.Variable
		matchesOutgoingPlaces = equals(api, transition.OutgoingPlaceCount, outgoingPlaceCount)
		tokenCountIncreasesOfTransition := tokenCountIncreasesByPlaceId.Lookup(transition.OutgoingPlaces[:]...)
		for j := range tokenCountIncreasesOfTransition {
			isPlaceIdInvalid := greaterThanOrEqual(api, j, transition.OutgoingPlaceCount)
			tokenCountIncreases := tokenCountIncreasesOfTransition[j]
			matchesOutgoingPlaces = api.And(matchesOutgoingPlaces, api.Or(isPlaceIdInvalid, tokenCountIncreases))
		}
		transitionFound = api.Or(transitionFound, api.And(matchesIncomingPlaces, matchesOutgoingPlaces))
	}

	tokenCountsDoNotChange := api.And(api.IsZero(incomingPlaceCount), api.IsZero(outgoingPlaceCount))
	api.AssertIsEqual(1, api.Or(transitionFound, tokenCountsDoNotChange))
}

func greaterThanOrEqual(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.Cmp(-1, api.Cmp(a, b))
}
