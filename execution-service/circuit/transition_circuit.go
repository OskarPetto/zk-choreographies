package circuit

import (
	"proof-service/domain"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/lookup/logderivlookup"
)

type TokenCountChanges struct {
	tokenCountDecreasesPerPlaceId *logderivlookup.Table
	tokenCountIncreasesPerPlaceId *logderivlookup.Table
	noChanges                     frontend.Variable
}

type TransitionCircuit struct {
	Model                 Model
	CurrentInstance       Instance
	NextInstance          Instance
	NextInstanceSignature Signature
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
	tokenCountChanges := circuit.compareTokenCounts(api)
	circuit.comparePublicKeys(api)
	participantId := findParticipantId(api, circuit.NextInstanceSignature, circuit.NextInstance)
	messageId := circuit.findMessageId(api)
	circuit.findTransition(api, tokenCountChanges, participantId, messageId)
	return nil
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

	// insert 1 at domain.MaxPlaceCount (default value of incomingPlaces and outgoingPlaces arrays)
	tokenCountDecreasesPerPlaceId.Insert(1)
	tokenCountIncreasesPerPlaceId.Insert(1)

	api.AssertIsLessOrEqual(tokenCountDecreasesCount, domain.MaxBranchingFactor)
	api.AssertIsLessOrEqual(tokenCountIncreasesCount, domain.MaxBranchingFactor)

	noChanges := api.And(api.IsZero(tokenCountDecreasesCount), api.IsZero(tokenCountIncreasesCount))

	return TokenCountChanges{
		tokenCountDecreasesPerPlaceId, tokenCountIncreasesPerPlaceId, noChanges,
	}
}

func (circuit *TransitionCircuit) comparePublicKeys(api frontend.API) {
	for i := range circuit.CurrentInstance.PublicKeys {
		currentPublicKey := circuit.CurrentInstance.PublicKeys[i]
		nextPublicKey := circuit.NextInstance.PublicKeys[i]
		api.AssertIsEqual(currentPublicKey.A.X, nextPublicKey.A.X)
		api.AssertIsEqual(currentPublicKey.A.Y, nextPublicKey.A.Y)
	}
}

func (circuit *TransitionCircuit) findMessageId(api frontend.API) frontend.Variable {
	var messageHashesAddedCount frontend.Variable = 0
	var addedMessageId frontend.Variable = domain.MaxMessageCount

	for messageId := range circuit.CurrentInstance.MessageHashes {
		currentMessageHash := circuit.CurrentInstance.MessageHashes[messageId]
		nextMessageHash := circuit.NextInstance.MessageHashes[messageId]
		var currentMessageHashIsDefault frontend.Variable = equals(api, currentMessageHash.Value, defaultMessageHash)
		var messageHashesMatch frontend.Variable = equals(api, currentMessageHash.Value, nextMessageHash.Value)
		api.AssertIsEqual(1, api.Or(currentMessageHashIsDefault, messageHashesMatch))

		messageHashAdded := api.IsZero(messageHashesMatch)
		addedMessageId = api.Select(messageHashAdded, messageId, addedMessageId)
		messageHashesAddedCount = api.Add(messageHashesAddedCount, messageHashAdded)
	}

	api.AssertIsLessOrEqual(messageHashesAddedCount, 1)

	return addedMessageId
}

func (circuit *TransitionCircuit) findTransition(api frontend.API, tokenCountChanges TokenCountChanges, participantId frontend.Variable, messageId frontend.Variable) {

	var transitionFound frontend.Variable = 0

	for _, transition := range circuit.Model.Transitions {
		participantMatches := api.Or(equals(api, transition.Participant, domain.MaxParticipantCount), equals(api, transition.Participant, participantId))
		messageMatches := equals(api, transition.Message, messageId)
		var tokenCountChangesMatch frontend.Variable = 1
		// returns 1 for default placeId (domain.MaxPlaceCount)
		incomingTokenCountsDecrease := tokenCountChanges.tokenCountDecreasesPerPlaceId.Lookup(transition.IncomingPlaces[:]...)
		// returns 1 for default placeId (domain.MaxPlaceCount)
		outgoingTokenCountsIncrease := tokenCountChanges.tokenCountIncreasesPerPlaceId.Lookup(transition.OutgoingPlaces[:]...)
		for j := range transition.IncomingPlaces {
			tokenCountChangesMatch = api.And(tokenCountChangesMatch, incomingTokenCountsDecrease[j])
			tokenCountChangesMatch = api.And(tokenCountChangesMatch, outgoingTokenCountsIncrease[j])
		}
		transitionMatches := api.And(api.And(tokenCountChangesMatch, participantMatches), messageMatches)
		transitionFound = api.Or(transitionFound, api.And(transition.IsValid, transitionMatches))
	}

	noMessageHashesAdded := equals(api, messageId, domain.MaxMessageCount)
	noChanges := api.And(tokenCountChanges.noChanges, noMessageHashesAdded)
	api.AssertIsEqual(1, api.Or(transitionFound, noChanges))
}
