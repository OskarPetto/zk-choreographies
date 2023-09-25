package circuit

import (
	"proof-service/domain"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/selector"
)

type TokenCountChanges struct {
	tokenCountDecreasesPerPlaceId [domain.MaxPlaceCount + 1]frontend.Variable
	tokenCountIncreasesPerPlaceId [domain.MaxPlaceCount + 1]frontend.Variable
	noChanges                     frontend.Variable
}

type TransitionCircuit struct {
	CurrentInstance       Instance
	NextInstance          Instance
	NextInstanceSignature Signature
	Model                 Model
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
	circuit.checkTransition(api, tokenCountChanges, participantId, messageId)
	return nil
}

func (circuit *TransitionCircuit) compareTokenCounts(api frontend.API) TokenCountChanges {
	var tokenCountDecreasesPerPlaceId [domain.MaxPlaceCount + 1]frontend.Variable
	var tokenCountIncreasesPerPlaceId [domain.MaxPlaceCount + 1]frontend.Variable
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

		tokenCountDecreasesPerPlaceId[placeId] = tokenCountDecreases
		tokenCountIncreasesPerPlaceId[placeId] = tokenCountIncreases

		api.AssertIsBoolean(nextTokenCount)
	}

	// insert 1 at domain.MaxPlaceCount (default value of incomingPlaces and outgoingPlaces arrays)
	tokenCountDecreasesPerPlaceId[domain.MaxPlaceCount] = 1
	tokenCountIncreasesPerPlaceId[domain.MaxPlaceCount] = 1

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
		api.AssertIsEqual(1, publicKeyEquals(api, currentPublicKey, nextPublicKey))
	}
}

func (circuit *TransitionCircuit) findMessageId(api frontend.API) frontend.Variable {
	var messageHashesAddedCount frontend.Variable = 0
	var sentMessageId frontend.Variable = domain.MaxMessageCount

	for messageId := range circuit.CurrentInstance.MessageHashes {
		currentMessageHash := circuit.CurrentInstance.MessageHashes[messageId]
		nextMessageHash := circuit.NextInstance.MessageHashes[messageId]
		var currentMessageHashIsZero frontend.Variable = 1
		var messageHashesMatch frontend.Variable = 1
		for i := range currentMessageHash.Value {
			currentMessageHashIsZero = api.And(currentMessageHashIsZero, api.IsZero(currentMessageHash.Value[i]))
			messageHashesMatch = api.And(messageHashesMatch, equals(api, currentMessageHash.Value[i], nextMessageHash.Value[i]))
		}
		api.AssertIsEqual(1, api.Or(currentMessageHashIsZero, messageHashesMatch))

		messageHashAdded := api.IsZero(messageHashesMatch)
		sentMessageId = api.Select(messageHashAdded, messageId, sentMessageId)
		messageHashesAddedCount = api.Add(messageHashesAddedCount, messageHashAdded)
	}

	api.AssertIsLessOrEqual(messageHashesAddedCount, 1)

	return sentMessageId
}

func (circuit *TransitionCircuit) checkTransition(api frontend.API, tokenCountChanges TokenCountChanges, participantId frontend.Variable, messageId frontend.Variable) {

	var transitionFound frontend.Variable = 0

	for _, transition := range circuit.Model.Transitions {
		participantMatches := api.Or(equals(api, transition.Participant, domain.MaxParticipantCount), equals(api, transition.Participant, participantId))
		messageMatches := equals(api, transition.Message, messageId)
		var tokenCountChangesMatch frontend.Variable = 1
		for j := range transition.IncomingPlaces {
			// returns 1 for default placeId (domain.MaxPlaceCount)
			incomingTokenCountDecreases := selector.Mux(api, transition.IncomingPlaces[j], tokenCountChanges.tokenCountDecreasesPerPlaceId[:]...)
			// returns 1 for default placeId (domain.MaxPlaceCount)
			outgoingTokenCountIncreases := selector.Mux(api, transition.OutgoingPlaces[j], tokenCountChanges.tokenCountIncreasesPerPlaceId[:]...)
			tokenCountChangesMatch = api.And(tokenCountChangesMatch, incomingTokenCountDecreases)
			tokenCountChangesMatch = api.And(tokenCountChangesMatch, outgoingTokenCountIncreases)
		}
		transitionMatches := api.And(api.And(tokenCountChangesMatch, participantMatches), messageMatches)
		transitionFound = api.Or(transitionFound, api.And(transition.IsValid, transitionMatches))
	}

	noMessageHashesAdded := equals(api, messageId, domain.MaxMessageCount)
	noChanges := api.And(tokenCountChanges.noChanges, noMessageHashesAdded)
	api.AssertIsEqual(1, api.Or(transitionFound, noChanges))
}
