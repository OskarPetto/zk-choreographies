package circuit

import (
	"execution-service/domain"
	"math/big"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/cmp"
	"github.com/consensys/gnark/std/selector"
)

type TokenCountChanges struct {
	placesWhereTokenCountDecreases [domain.MaxBranchingFactor]frontend.Variable
	placesWhereTokenCountIncreases [domain.MaxBranchingFactor]frontend.Variable
	noChanges                      frontend.Variable
}

type ConstraintMessageIds struct {
	MessageIds [domain.MaxConstraintMessageCount]frontend.Variable
}

type TransitionCircuit struct {
	Model                   Model
	CurrentInstance         Instance
	NextInstance            Instance
	Transition              Transition
	SenderAuthentication    Authentication
	RecipientAuthentication Authentication
	ConstraintInput         ConstraintInput
}

func NewTransitionCircuit() TransitionCircuit {
	return TransitionCircuit{
		SenderAuthentication: Authentication{
			MerkleProof: MerkleProof{
				MerkleProof: merkle.MerkleProof{
					Path: make([]frontend.Variable, domain.MaxParticipantDepth+1),
				},
			},
		},
		RecipientAuthentication: Authentication{
			MerkleProof: MerkleProof{
				MerkleProof: merkle.MerkleProof{
					Path: make([]frontend.Variable, domain.MaxParticipantDepth+1),
				},
			},
		},
		Transition: Transition{
			MerkleProof: MerkleProof{
				MerkleProof: merkle.MerkleProof{
					Path: make([]frontend.Variable, domain.MaxTransitionDepth+1),
				},
			},
		},
	}
}

func (circuit *TransitionCircuit) Define(api frontend.API) error {
	err := checkModelHash(api, circuit.Model, circuit.CurrentInstance)
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
	err = checkAuthentication(api, circuit.SenderAuthentication, circuit.NextInstance)
	if err != nil {
		return err
	}
	err = checkAuthentication(api, circuit.RecipientAuthentication, circuit.NextInstance)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.CurrentInstance.Model, circuit.NextInstance.Model)
	circuit.comparePublicKeys(api)
	tokenCountChanges := circuit.compareTokenCounts(api)
	addedMessageId := circuit.findAddedMessageId(api)
	constraintInputMessageIds, err := circuit.findConstraintInputMessageIds(api)
	if err != nil {
		return err
	}
	return circuit.checkTransition(api, tokenCountChanges, addedMessageId, constraintInputMessageIds)
}

func (circuit *TransitionCircuit) findConstraintInputMessageIds(api frontend.API) (ConstraintMessageIds, error) {
	var messageIds [domain.MaxConstraintMessageCount]frontend.Variable
	var messageHashesMatchCount frontend.Variable = 0

	for i := 0; i < domain.MaxConstraintMessageCount; i++ {
		messageIds[i] = domain.EmptyMessageId
	}
	for i, integerMessage := range circuit.ConstraintInput.IntegerMessages {
		salt := circuit.ConstraintInput.Salts[i]
		mimc, err := mimc.NewMiMC(api)
		if err != nil {
			return ConstraintMessageIds{}, err
		}
		mimc.Write(integerMessage)
		mimc.Write(salt)
		messageHash := mimc.Sum()

		for messageId, messageHashForMessageId := range circuit.CurrentInstance.MessageHashes {
			messageHashesMatch := equals(api, messageHash, messageHashForMessageId)
			for i := 0; i < domain.MaxConstraintMessageCount; i++ {
				isCorrectIndex := equals(api, messageHashesMatchCount, i)
				shouldWrite := api.And(isCorrectIndex, messageHashesMatch)
				messageIds[i] = api.Select(shouldWrite, messageId, messageIds[i])
			}
			messageHashesMatchCount = api.Add(messageHashesMatchCount, messageHashesMatch)
		}
	}
	return ConstraintMessageIds{
		MessageIds: messageIds,
	}, nil
}

func (circuit *TransitionCircuit) comparePublicKeys(api frontend.API) {
	api.AssertIsEqual(circuit.CurrentInstance.PublicKeyRoot, circuit.NextInstance.PublicKeyRoot)
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

func (circuit *TransitionCircuit) findAddedMessageId(api frontend.API) frontend.Variable {
	var messageHashesAddedCount frontend.Variable = 0
	var addedMessageId frontend.Variable = domain.EmptyMessageId

	for messageId := range circuit.CurrentInstance.MessageHashes {
		currentMessageHash := circuit.CurrentInstance.MessageHashes[messageId]
		nextMessageHash := circuit.NextInstance.MessageHashes[messageId]

		messageHashesMatch := equals(api, currentMessageHash, nextMessageHash)
		messageHashAdded := api.IsZero(messageHashesMatch)
		addedMessageId = api.Select(messageHashAdded, messageId, addedMessageId)
		messageHashesAddedCount = api.Add(messageHashesAddedCount, messageHashAdded)
	}

	api.AssertIsLessOrEqual(messageHashesAddedCount, 1)

	return addedMessageId
}

func (circuit *TransitionCircuit) checkTransition(api frontend.API, tokenCountChanges TokenCountChanges, addedMessageId frontend.Variable, constraintMessageIds ConstraintMessageIds) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	circuit.Transition.MerkleProof.CheckRootHash(api, circuit.Model.TransitionRoot)
	circuit.Transition.MerkleProof.VerifyProof(api, mimc)
	checkTransitionHash(api, circuit.Transition.MerkleProof.MerkleProof.Path[0], circuit.Transition)

	transition := circuit.Transition
	senderParticipantId := circuit.SenderAuthentication.MerkleProof.Index
	recipientParticipantId := circuit.RecipientAuthentication.MerkleProof.Index

	var tokenCountChangesMatch frontend.Variable = 1
	for _, incomingPlace := range transition.IncomingPlaces {
		var tokenCountDecreasesAtIncomingPlace frontend.Variable = 0
		for _, placeWhereTokenCountDecreases := range tokenCountChanges.placesWhereTokenCountDecreases {
			tokenCountDecreasesAtIncomingPlace = api.Or(tokenCountDecreasesAtIncomingPlace, equals(api, incomingPlace, placeWhereTokenCountDecreases))
		}
		tokenCountChangesMatch = api.And(tokenCountChangesMatch, tokenCountDecreasesAtIncomingPlace)
	}
	for _, outgoingPlace := range transition.OutgoingPlaces {
		var tokenCountIncreasesAtOutgoingPlace frontend.Variable = 0
		for _, placeWhereTokenCountIncreases := range tokenCountChanges.placesWhereTokenCountIncreases {
			tokenCountIncreasesAtOutgoingPlace = api.Or(tokenCountIncreasesAtOutgoingPlace, equals(api, outgoingPlace, placeWhereTokenCountIncreases))
		}
		tokenCountChangesMatch = api.And(tokenCountChangesMatch, tokenCountIncreasesAtOutgoingPlace)
	}

	senderMatches := api.Or(equals(api, transition.Sender, domain.EmptyParticipantId), equals(api, transition.Sender, senderParticipantId))
	recipientMatches := api.Or(equals(api, transition.Recipient, domain.EmptyParticipantId), equals(api, transition.Recipient, recipientParticipantId))
	messageMatches := equals(api, transition.Message, addedMessageId)
	constraintSatisfied := evaluateConstraint(api, transition.Constraint, circuit.ConstraintInput, constraintMessageIds)

	transitionMatches := api.And(senderMatches, recipientMatches)
	transitionMatches = api.And(transitionMatches, tokenCountChangesMatch)
	transitionMatches = api.And(transitionMatches, messageMatches)
	transitionMatches = api.And(transitionMatches, constraintSatisfied)

	noMessageHashesAdded := equals(api, addedMessageId, domain.EmptyMessageId)
	noChanges := api.And(tokenCountChanges.noChanges, noMessageHashesAdded)
	api.AssertIsEqual(1, api.Or(transitionMatches, noChanges))
	return nil
}

func evaluateConstraint(api frontend.API, constraint Constraint, input ConstraintInput, constraintMessageIds ConstraintMessageIds) frontend.Variable {
	lhs := constraint.Offset
	var allMessageIdsMatch frontend.Variable = 1
	for i, messageId := range constraint.MessageIds {
		coefficient := constraint.Coefficients[i]
		expectedMessageId := constraintMessageIds.MessageIds[i]

		coefficientIsZero := api.IsZero(coefficient)
		messageIdsMatch := equals(api, expectedMessageId, messageId)

		allMessageIdsMatch = api.And(allMessageIdsMatch, api.Or(coefficientIsZero, messageIdsMatch))
		integerMessage := input.IntegerMessages[i]
		lhs = api.MulAcc(lhs, coefficient, integerMessage)
	}

	comparator := cmp.NewBoundedComparator(api, big.NewInt(1<<32-1), false)
	var comparisons [5]frontend.Variable
	comparisons[0] = equals(api, lhs, 0)                                // equal
	comparisons[1] = comparator.IsLess(0, lhs)                          // gt
	comparisons[2] = api.IsZero(api.Or(comparisons[0], comparisons[1])) // lt
	comparisons[3] = api.Or(comparisons[0], comparisons[1])             // gte
	comparisons[4] = api.Or(comparisons[0], comparisons[2])             // lte
	result := selector.Mux(api, constraint.ComparisonOperator, comparisons[:]...)
	return api.And(allMessageIdsMatch, result)
}

func checkTransitionHash(api frontend.API, hash frontend.Variable, transition Transition) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(transition.IncomingPlaces[:]...)
	mimc.Write(transition.OutgoingPlaces[:]...)
	mimc.Write(transition.Sender)
	mimc.Write(transition.Recipient)
	mimc.Write(transition.Message)
	mimc.Write(transition.Constraint.Coefficients[:]...)
	mimc.Write(transition.Constraint.MessageIds[:]...)
	mimc.Write(transition.Constraint.Offset)
	mimc.Write(transition.Constraint.ComparisonOperator)
	api.AssertIsEqual(hash, mimc.Sum())
	return nil
}
