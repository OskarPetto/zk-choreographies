package circuit

import (
	"execution-service/domain"

	"github.com/consensys/gnark/frontend"
)

type TokenCountChanges struct {
	placesWhereTokenCountDecreases [domain.MaxBranchingFactor]frontend.Variable
	placesWhereTokenCountIncreases [domain.MaxBranchingFactor]frontend.Variable
	noChanges                      frontend.Variable
}

type TransitionCircuit struct {
	Model           Model
	CurrentInstance Instance
	NextInstance    Instance
	NextSignature   Signature
}

func (circuit *TransitionCircuit) Define(api frontend.API) error {
	err := checkModelHash(api, circuit.Model)
	if err != nil {
		return err
	}
	err = checkInstanceHash(api, circuit.CurrentInstance)
	if err != nil {
		return err
	}
	err = checkInstanceHash(api, circuit.NextInstance)
	if err != nil {
		return err
	}
	err = checkSignature(api, circuit.NextSignature, circuit.NextInstance)
	if err != nil {
		return err
	}
	api.AssertIsLessOrEqual(circuit.CurrentInstance.CreatedAt, circuit.NextInstance.CreatedAt)
	tokenCountChanges := circuit.compareTokenCounts(api)
	circuit.comparePublicKeys(api)
	participantId := findParticipantId(api, circuit.NextSignature, circuit.NextInstance)
	messageId := circuit.findMessageId(api)
	api.Println(tokenCountChanges.placesWhereTokenCountDecreases[:]...)
	api.Println(tokenCountChanges.placesWhereTokenCountIncreases[:]...)
	circuit.findTransition(api, tokenCountChanges, participantId, messageId)
	return nil
}

func (circuit *TransitionCircuit) compareTokenCounts(api frontend.API) TokenCountChanges {
	var placesWhereTokenCountDecreases [domain.MaxBranchingFactor]frontend.Variable
	var placesWhereTokenCountIncreases [domain.MaxBranchingFactor]frontend.Variable
	for i := 0; i < domain.MaxBranchingFactor; i++ {
		placesWhereTokenCountDecreases[i] = domain.OutOfBoundsPlaceId
		placesWhereTokenCountIncreases[i] = domain.OutOfBoundsPlaceId
	}
	var tokenCountDecreasesCount frontend.Variable = 0
	var tokenCountIncreasesCount frontend.Variable = 0

	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId]
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId]

		isTokenCount := api.Or(equals(api, nextTokenCount, 0), equals(api, nextTokenCount, 1))
		tokenChange := api.Sub(nextTokenCount, currentTokenCount)
		tokenCountStaysTheSame := api.IsZero(tokenChange)
		api.AssertIsEqual(1, api.Or(isTokenCount, tokenCountStaysTheSame))

		tokenCountDecreases := equals(api, tokenChange, -1)
		tokenCountIncreases := equals(api, tokenChange, 1)
		api.AssertIsEqual(1, api.Or(api.Or(tokenCountStaysTheSame, tokenCountDecreases), tokenCountIncreases))

		for i := 0; i < domain.MaxBranchingFactor; i++ {
			isCorrectIndex := equals(api, tokenCountDecreasesCount, i)
			shouldWrite := api.And(tokenCountDecreases, isCorrectIndex)
			placesWhereTokenCountDecreases[i] = api.Select(shouldWrite, placeId, placesWhereTokenCountDecreases[i])
		}
		for i := 0; i < domain.MaxBranchingFactor; i++ {
			isCorrectIndex := equals(api, tokenCountIncreasesCount, i)
			shouldWrite := api.And(tokenCountIncreases, isCorrectIndex)
			placesWhereTokenCountIncreases[i] = api.Select(shouldWrite, placeId, placesWhereTokenCountIncreases[i])
		}

		tokenCountDecreasesCount = api.Add(tokenCountDecreasesCount, tokenCountDecreases)
		tokenCountIncreasesCount = api.Add(tokenCountIncreasesCount, tokenCountIncreases)
	}

	api.AssertIsLessOrEqual(tokenCountDecreasesCount, domain.MaxBranchingFactor)
	api.AssertIsLessOrEqual(tokenCountIncreasesCount, domain.MaxBranchingFactor)

	noChanges := api.And(api.IsZero(tokenCountDecreasesCount), api.IsZero(tokenCountIncreasesCount))

	return TokenCountChanges{
		placesWhereTokenCountDecreases, placesWhereTokenCountIncreases, noChanges,
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
	var addedMessageId frontend.Variable = domain.EmptyMessageId

	for messageId := range circuit.CurrentInstance.MessageHashes {
		currentMessageHash := circuit.CurrentInstance.MessageHashes[messageId]
		nextMessageHash := circuit.NextInstance.MessageHashes[messageId]
		var wasMessageHashEmpty frontend.Variable = equals(api, currentMessageHash, emptyMessageHash)
		var messageHashesMatch frontend.Variable = equals(api, currentMessageHash, nextMessageHash)
		api.AssertIsEqual(1, api.Or(wasMessageHashEmpty, messageHashesMatch))

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
		participantMatches := api.Or(equals(api, transition.Participant, domain.EmptyParticipantId), equals(api, transition.Participant, participantId))
		messageMatches := equals(api, transition.Message, messageId)
		var tokenCountChangesMatch frontend.Variable = 1
		for _, incomingPlace := range transition.IncomingPlaces {
			var tokenCountDecreases frontend.Variable = 0
			for _, placeWhereTokenCountDecreases := range tokenCountChanges.placesWhereTokenCountDecreases {
				tokenCountDecreases = api.Or(tokenCountDecreases, equals(api, incomingPlace, placeWhereTokenCountDecreases))
			}
			tokenCountChangesMatch = api.And(tokenCountChangesMatch, tokenCountDecreases)
		}
		for _, outgoingPlace := range transition.OutgoingPlaces {
			var tokenCountIncreases frontend.Variable = 0
			for _, placeWhereTokenCountIncreases := range tokenCountChanges.placesWhereTokenCountIncreases {
				tokenCountIncreases = api.Or(tokenCountIncreases, equals(api, outgoingPlace, placeWhereTokenCountIncreases))
			}
			tokenCountChangesMatch = api.And(tokenCountChangesMatch, tokenCountIncreases)
		}
		transitionMatches := api.And(api.And(tokenCountChangesMatch, participantMatches), messageMatches)
		transitionFound = api.Or(transitionFound, api.And(transition.IsValid, transitionMatches))
	}

	noMessageHashesAdded := equals(api, messageId, domain.MaxMessageCount)
	noChanges := api.And(tokenCountChanges.noChanges, noMessageHashesAdded)
	api.AssertIsEqual(1, api.Or(transitionFound, noChanges))
}
