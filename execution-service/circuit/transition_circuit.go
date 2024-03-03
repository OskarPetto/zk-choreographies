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

type IntermediateResult struct {
	matchesTransition frontend.Variable
	noChanges         frontend.Variable
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
	err = circuit.checkTransition(api)
	if err != nil {
		return err
	}
	tokenCountChanges := circuit.checkTokenCounts(api)
	addedMessageHashes := circuit.checkAddedMessageHashes(api)
	participantsMatch := circuit.checkParticipants(api)
	constrainInputsMatch, err := circuit.checkConstraintInput(api)
	if err != nil {
		return err
	}
	constraintSatisfied := evaluateConstraint(api, circuit.Transition.Constraint, circuit.ConstraintInput)

	transitionMatches := api.And(participantsMatch, tokenCountChanges.matchesTransition)
	transitionMatches = api.And(transitionMatches, addedMessageHashes.matchesTransition)
	transitionMatches = api.And(transitionMatches, constrainInputsMatch)
	transitionMatches = api.And(transitionMatches, constraintSatisfied)

	noChanges := api.And(tokenCountChanges.noChanges, addedMessageHashes.noChanges)

	api.AssertIsEqual(1, api.Or(transitionMatches, noChanges))
	return nil
}

func (circuit *TransitionCircuit) checkTransition(api frontend.API) error {
	circuit.Transition.MerkleProof.CheckRootHash(api, circuit.Model.TransitionRoot)
	err := circuit.Transition.MerkleProof.VerifyProof(api)
	if err != nil {
		return err
	}
	err = checkTransitionHash(api, circuit.Transition.MerkleProof.MerkleProof.Path[0], circuit.Transition)
	if err != nil {
		return err
	}
	return nil
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

func (circuit *TransitionCircuit) checkConstraintInput(api frontend.API) (frontend.Variable, error) {
	var constraintInputMatchesTransition frontend.Variable = 1
	for i, referencedMessageId := range circuit.Transition.Constraint.MessageIds {
		noMessageReferenced := equals(api, referencedMessageId, domain.EmptyMessageId)
		var messageHash frontend.Variable = emptyMessageHash
		for messageId, messageHashForMessageId := range circuit.CurrentInstance.MessageHashes {
			messageIdsMatch := equals(api, messageId, referencedMessageId)
			messageHash = api.Select(messageIdsMatch, messageHashForMessageId, messageHash)
		}

		providedMessage := circuit.ConstraintInput.Messages[i]
		mimc, err := mimc.NewMiMC(api)
		if err != nil {
			return 0, err
		}
		mimc.Write(providedMessage.IntegerMessage)
		mimc.Write(providedMessage.Instance)
		mimc.Write(providedMessage.Salt)
		providedMessageHash := mimc.Sum()
		messageHashesMatch := equals(api, messageHash, providedMessageHash)

		constraintInputMatchesTransition = api.And(constraintInputMatchesTransition, api.Or(noMessageReferenced, messageHashesMatch))
	}
	return constraintInputMatchesTransition, nil
}

func (circuit *TransitionCircuit) checkTokenCounts(api frontend.API) IntermediateResult {
	var noChanges frontend.Variable = 1
	var tokenCountChangesMatchTransition frontend.Variable = 1

	for placeId := range circuit.CurrentInstance.TokenCounts {
		currentTokenCount := circuit.CurrentInstance.TokenCounts[placeId]
		nextTokenCount := circuit.NextInstance.TokenCounts[placeId]

		isTokenCount := api.Or(equals(api, nextTokenCount, 0), equals(api, nextTokenCount, 1))
		tokenChange := api.Sub(nextTokenCount, currentTokenCount)
		tokenCountStaysTheSame := api.IsZero(tokenChange)
		api.AssertIsEqual(1, api.Or(isTokenCount, tokenCountStaysTheSame))

		tokenCountDecreases := equals(api, tokenChange, -1)
		tokenCountIncreases := equals(api, tokenChange, 1)

		var isIncomingPlace frontend.Variable = 0
		for _, incomingPlace := range circuit.Transition.IncomingPlaces {
			isIncomingPlace = api.Or(isIncomingPlace, equals(api, placeId, incomingPlace))
		}
		isNotIncomingPlace := api.IsZero(isIncomingPlace)

		var isOutgoingPlace frontend.Variable = 0
		for _, outgoingPlace := range circuit.Transition.OutgoingPlaces {
			isOutgoingPlace = api.Or(isOutgoingPlace, equals(api, placeId, outgoingPlace))
		}
		isNotOutgoingPlace := api.IsZero(isOutgoingPlace)

		tokenCountDecreasesAtIncomingPlace := api.Or(isNotIncomingPlace, tokenCountDecreases)
		tokenCountIncreasesAtOutgoingPlace := api.Or(isNotOutgoingPlace, tokenCountIncreases)
		tokenCountStaysTheSameAtOtherPlace := api.Or(api.Or(isIncomingPlace, isOutgoingPlace), tokenCountStaysTheSame)

		tokenCountChangesMatchTransition = api.And(tokenCountChangesMatchTransition, tokenCountDecreasesAtIncomingPlace)
		tokenCountChangesMatchTransition = api.And(tokenCountChangesMatchTransition, tokenCountIncreasesAtOutgoingPlace)
		tokenCountChangesMatchTransition = api.And(tokenCountChangesMatchTransition, tokenCountStaysTheSameAtOtherPlace)

		noChanges = api.And(noChanges, tokenCountStaysTheSame)
	}

	return IntermediateResult{
		matchesTransition: tokenCountChangesMatchTransition,
		noChanges:         noChanges,
	}
}

func (circuit *TransitionCircuit) checkAddedMessageHashes(api frontend.API) IntermediateResult {
	var noChanges frontend.Variable = 1
	var addedMessageHashesMatchTransition frontend.Variable = 1

	for messageId := range circuit.CurrentInstance.MessageHashes {
		currentMessageHash := circuit.CurrentInstance.MessageHashes[messageId]
		nextMessageHash := circuit.NextInstance.MessageHashes[messageId]

		messageHashesMatch := equals(api, currentMessageHash, nextMessageHash)
		messageHashAdded := api.IsZero(messageHashesMatch)

		isInitiatingMessage := equals(api, messageId, circuit.Transition.InitiatingMessage)
		isRespondingMessage := equals(api, messageId, circuit.Transition.RespondingMessage)

		messageHashAddedAtInitiatingMessage := api.Or(api.IsZero(isInitiatingMessage), messageHashAdded)
		messageHashAddedAtRespondingMessage := api.Or(api.IsZero(isRespondingMessage), messageHashAdded)
		messageHashesMatchAtOtherMessage := api.Or(api.Or(isInitiatingMessage, isRespondingMessage), messageHashesMatch)

		addedMessageHashesMatchTransition = api.And(addedMessageHashesMatchTransition, messageHashAddedAtInitiatingMessage)
		addedMessageHashesMatchTransition = api.And(addedMessageHashesMatchTransition, messageHashAddedAtRespondingMessage)
		addedMessageHashesMatchTransition = api.And(addedMessageHashesMatchTransition, messageHashesMatchAtOtherMessage)

		noChanges = api.And(noChanges, messageHashesMatch)
	}

	return IntermediateResult{
		matchesTransition: addedMessageHashesMatchTransition,
		noChanges:         noChanges,
	}
}

func (circuit *TransitionCircuit) checkParticipants(api frontend.API) frontend.Variable {
	transition := circuit.Transition
	senderParticipantId := circuit.InitiatingParticipantAuthentication.MerkleProof.Index
	recipientParticipantId := circuit.RespondingParticipantAuthentication.MerkleProof.Index
	initiatingParticipantMatches := api.Or(equals(api, transition.InitiatingParticipant, domain.EmptyParticipantId), equals(api, transition.InitiatingParticipant, senderParticipantId))
	respondingParticipantMatches := api.Or(equals(api, transition.RespondingParticipant, domain.EmptyParticipantId), equals(api, transition.RespondingParticipant, recipientParticipantId))
	return api.And(initiatingParticipantMatches, respondingParticipantMatches)
}

func evaluateConstraint(api frontend.API, constraint Constraint, input ConstraintInput) frontend.Variable {
	lhs := constraint.Offset
	for i := range constraint.MessageIds {
		coefficient := constraint.Coefficients[i]
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
	return result
}
