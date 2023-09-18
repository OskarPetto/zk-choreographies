package circuit

import (
	"proof-service/domain"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/selector"
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
	uapi.ByteAssertEq(circuit.NextInstance.PlaceCount, circuit.PetriNet.PlaceCount)
	circuit.checkTokenCounts(api)
	return nil
}

func (circuit *ExecutionCircuit) checkTokenCounts(api frontend.API) {
	var tokenCountDecreases [domain.MaxPlaceCount]frontend.Variable
	var tokenCountIncreases [domain.MaxPlaceCount]frontend.Variable
	var incomingPlaceCount frontend.Variable
	var outgoingPlaceCount frontend.Variable
	incomingPlaceCount = 0
	outgoingPlaceCount = 0
	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId].Val
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId].Val

		tokenCountStaysTheSame := equals(api, nextTokenCount, currentTokenCount)
		doesTokenCountDecrease := equals(api, nextTokenCount, api.Sub(currentTokenCount, 1))
		doesTokenCountIncrease := equals(api, nextTokenCount, api.Add(currentTokenCount, 1))
		api.AssertIsEqual(1, api.Or(api.Or(tokenCountStaysTheSame, doesTokenCountDecrease), doesTokenCountIncrease))

		incomingPlaceCount = api.Add(incomingPlaceCount, doesTokenCountDecrease)
		outgoingPlaceCount = api.Add(outgoingPlaceCount, doesTokenCountIncrease)

		tokenCountDecreases[placeId] = doesTokenCountDecrease
		tokenCountIncreases[placeId] = doesTokenCountIncrease

		api.AssertIsBoolean(nextTokenCount)
	}

	api.AssertIsLessOrEqual(incomingPlaceCount, api.Sub(domain.MaxBranchingFactor, 1))
	api.AssertIsLessOrEqual(outgoingPlaceCount, api.Sub(domain.MaxBranchingFactor, 1))

	var transitionFound frontend.Variable
	transitionFound = 0

	for i := range circuit.PetriNet.Transitions {
		transition := circuit.PetriNet.Transitions[i]
		var matchesIncomingPlaces frontend.Variable
		matchesIncomingPlaces = equals(api, transition.IncomingPlaceCount.Val, incomingPlaceCount)
		for j := range transition.IncomingPlaces {
			isPlaceIdInvalid := greaterThanOrEqual(api, j, transition.IncomingPlaceCount.Val)
			incomingPlaceId := transition.IncomingPlaces[j]
			doesTokenCountDecrease := selector.Mux(api, incomingPlaceId.Val, tokenCountDecreases[:]...)
			matchesIncomingPlaces = api.And(matchesIncomingPlaces, api.Or(isPlaceIdInvalid, doesTokenCountDecrease))
		}
		var matchesOutgoingPlaces frontend.Variable
		matchesOutgoingPlaces = equals(api, transition.OutgoingPlaceCount.Val, outgoingPlaceCount)
		for j := range transition.OutgoingPlaces {
			isPlaceIdInvalid := greaterThanOrEqual(api, j, transition.OutgoingPlaceCount.Val)
			outgoingPlaceId := transition.OutgoingPlaces[j]
			doesTokenCountIncrease := selector.Mux(api, outgoingPlaceId.Val, tokenCountIncreases[:]...)
			matchesOutgoingPlaces = api.And(matchesOutgoingPlaces, api.Or(isPlaceIdInvalid, doesTokenCountIncrease))
		}
		transitionFound = api.Or(transitionFound, api.And(matchesIncomingPlaces, matchesOutgoingPlaces))
	}

	tokenCountsDoNotChange := api.And(api.IsZero(incomingPlaceCount), api.IsZero(outgoingPlaceCount))
	api.AssertIsEqual(1, api.Or(transitionFound, tokenCountsDoNotChange))
}

func greaterThanOrEqual(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.Cmp(-1, api.Cmp(a, b))
}
