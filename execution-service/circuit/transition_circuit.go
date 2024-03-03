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

type AddedMessageIds struct {
	messageIds      [2]frontend.Variable
	noMessagesAdded frontend.Variable
}

type TokenCountChanges struct {
	placesWhereTokenCountDecreases [domain.MaxBranchingFactor]frontend.Variable
	placesWhereTokenCountIncreases [domain.MaxBranchingFactor]frontend.Variable
	noChanges                      frontend.Variable
}

type ConstraintMessageIds struct {
	MessageIds [domain.MaxMessageCountInConstraints]frontend.Variable
}

type TransitionCircuit struct {
	Model                               Model
	CurrentInstance                     Instance
	NextInstance                        Instance
	Transition                          Transition
	InitiatingParticipantAuthentication Authentication
	RespondingParticipantAuthentication Authentication
	ConstraintInput                     ConstraintInput
}

func NewTransitionCircuit() TransitionCircuit {
	return TransitionCircuit{
		InitiatingParticipantAuthentication: Authentication{
			MerkleProof: MerkleProof{
				MerkleProof: merkle.MerkleProof{
					Path: make([]frontend.Variable, domain.MaxParticipantDepth+1),
				},
			},
		},
		RespondingParticipantAuthentication: Authentication{
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
	api.AssertIsEqual(circuit.CurrentInstance.Model, circuit.NextInstance.Model)
	err = checkAuthentication(api, circuit.InitiatingParticipantAuthentication, circuit.NextInstance)
	if err != nil {
		return err
	}
	err = checkAuthentication(api, circuit.RespondingParticipantAuthentication, circuit.NextInstance)
	if err != nil {
		return err
	}
	api.AssertIsEqual(circuit.CurrentInstance.PublicKeyRoot, circuit.NextInstance.PublicKeyRoot)
	tokenCountChanges := circuit.compareTokenCounts(api)
	addedMessageIds := circuit.findAddedMessageIds(api)
	constraintInputMessageIds, err := circuit.findConstraintInputMessageIds(api)
	if err != nil {
		return err
	}
	return circuit.checkTransition(api, tokenCountChanges, addedMessageIds, constraintInputMessageIds)
}

func (circuit *TransitionCircuit) findConstraintInputMessageIds(api frontend.API) (ConstraintMessageIds, error) {
	var messageIds [domain.MaxMessageCountInConstraints]frontend.Variable
	var messageHashesMatchCount frontend.Variable = 0

	for i := 0; i < domain.MaxMessageCountInConstraints; i++ {
		messageIds[i] = domain.EmptyMessageId
	}
	for _, message := range circuit.ConstraintInput.Messages {
		mimc, err := mimc.NewMiMC(api)
		if err != nil {
			return ConstraintMessageIds{}, err
		}
		mimc.Write(message.IntegerMessage)
		mimc.Write(message.Instance)
		mimc.Write(message.Salt)
		messageHash := mimc.Sum()

		for messageId, messageHashForMessageId := range circuit.CurrentInstance.MessageHashes {
			messageHashesMatch := equals(api, messageHash, messageHashForMessageId)
			for i := 0; i < domain.MaxMessageCountInConstraints; i++ {
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

func (circuit *TransitionCircuit) findAddedMessageIds(api frontend.API) AddedMessageIds {
	var messageHashesAddedCount frontend.Variable = 0
	addedMessageIds := [2]frontend.Variable{domain.EmptyMessageId, domain.EmptyMessageId}

	for messageId := range circuit.CurrentInstance.MessageHashes {
		currentMessageHash := circuit.CurrentInstance.MessageHashes[messageId]
		nextMessageHash := circuit.NextInstance.MessageHashes[messageId]

		messageHashesMatch := equals(api, currentMessageHash, nextMessageHash)
		messageHashAdded := api.IsZero(messageHashesMatch)

		for i := 0; i < 2; i++ {
			isCorrectIndex := equals(api, messageHashesAddedCount, i)
			shouldWrite := api.And(messageHashAdded, isCorrectIndex)
			addedMessageIds[i] = api.Select(shouldWrite, messageId, addedMessageIds[i])
		}

		messageHashesAddedCount = api.Add(messageHashesAddedCount, messageHashAdded)
	}

	api.AssertIsLessOrEqual(messageHashesAddedCount, 2)

	return AddedMessageIds{
		messageIds:      addedMessageIds,
		noMessagesAdded: api.IsZero(messageHashesAddedCount),
	}
}

func (circuit *TransitionCircuit) checkTransition(api frontend.API, tokenCountChanges TokenCountChanges, addedMessageIds AddedMessageIds, constraintMessageIds ConstraintMessageIds) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	circuit.Transition.MerkleProof.CheckRootHash(api, circuit.Model.TransitionRoot)
	circuit.Transition.MerkleProof.VerifyProof(api, mimc)
	checkTransitionHash(api, circuit.Transition.MerkleProof.MerkleProof.Path[0], circuit.Transition)

	transition := circuit.Transition
	senderParticipantId := circuit.InitiatingParticipantAuthentication.MerkleProof.Index
	recipientParticipantId := circuit.RespondingParticipantAuthentication.MerkleProof.Index

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

	initiatingParticipantMatches := api.Or(equals(api, transition.InitiatingParticipant, domain.EmptyParticipantId), equals(api, transition.InitiatingParticipant, senderParticipantId))
	respondingParticipantMatches := api.Or(equals(api, transition.RespondingParticipant, domain.EmptyParticipantId), equals(api, transition.RespondingParticipant, recipientParticipantId))
	messageMatches := addedMessagesMatch(api, transition.InitiatingMessage, transition.RespondingMessage, addedMessageIds)
	constraintSatisfied := evaluateConstraint(api, transition.Constraint, circuit.ConstraintInput, constraintMessageIds)

	transitionMatches := api.And(initiatingParticipantMatches, respondingParticipantMatches)
	transitionMatches = api.And(transitionMatches, tokenCountChangesMatch)
	transitionMatches = api.And(transitionMatches, messageMatches)
	transitionMatches = api.And(transitionMatches, constraintSatisfied)

	noChanges := api.And(tokenCountChanges.noChanges, addedMessageIds.noMessagesAdded)

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

		message := input.Messages[i]
		lhs = api.MulAcc(lhs, coefficient, message.IntegerMessage)
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
	mimc.Write(transition.InitiatingParticipant)
	mimc.Write(transition.RespondingParticipant)
	mimc.Write(transition.InitiatingMessage)
	mimc.Write(transition.Constraint.Coefficients[:]...)
	mimc.Write(transition.Constraint.MessageIds[:]...)
	mimc.Write(transition.Constraint.Offset)
	mimc.Write(transition.Constraint.ComparisonOperator)
	api.AssertIsEqual(hash, mimc.Sum())
	return nil
}

func addedMessagesMatch(api frontend.API, initiatingMessage frontend.Variable, respondingMessage frontend.Variable, addedMessageIds AddedMessageIds) frontend.Variable {
	return api.Or(
		api.And(equals(api, initiatingMessage, addedMessageIds.messageIds[0]), equals(api, respondingMessage, addedMessageIds.messageIds[1])),
		api.And(equals(api, initiatingMessage, addedMessageIds.messageIds[1]), equals(api, respondingMessage, addedMessageIds.messageIds[0])),
	)
}
