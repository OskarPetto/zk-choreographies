package circuit_test

import (
	"execution-service/circuit"
	"execution-service/domain"
	"execution-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit = circuit.NewTransitionCircuit()
var transitionStates = testdata.GetModel2States(signatureParameters)

func TestTransition_NoChange(t *testing.T) {
	model := transitionStates[4].Model
	currentInstance := transitionStates[4].Instance
	transition := transitionStates[4].Transition
	signature := transitionStates[4].InitiatingParticipantSignature
	constraintInput := transitionStates[4].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(currentInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(currentInstance, signature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(currentInstance, signature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransition_Transitions(t *testing.T) {
	for i := 0; i < len(transitionStates)-1; i++ {
		model := transitionStates[i].Model
		currentInstance := transitionStates[i].Instance
		nextInstance := transitionStates[i+1].Instance
		transition := transitionStates[i+1].Transition
		initiatingParticipantSignature := transitionStates[i+1].InitiatingParticipantSignature
		respondingParticipantSignature := transitionStates[i+1].RespondingParticipantSignature
		constraintInput := transitionStates[i+1].ConstraintInput

		initiatingParticipantAuthentication := circuit.ToAuthentication(nextInstance, initiatingParticipantSignature)
		respondingParticipantAuthentication := initiatingParticipantAuthentication
		if respondingParticipantSignature != nil {
			respondingParticipantAuthentication = circuit.ToAuthentication(nextInstance, *respondingParticipantSignature)
		}
		witness := circuit.TransitionCircuit{
			Model:                               circuit.FromModel(model),
			CurrentInstance:                     circuit.FromInstance(currentInstance),
			NextInstance:                        circuit.FromInstance(nextInstance),
			Transition:                          circuit.ToTransition(model, transition),
			InitiatingParticipantAuthentication: initiatingParticipantAuthentication,
			RespondingParticipantAuthentication: respondingParticipantAuthentication,
			ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
		}

		err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestTransition_InvalidModelHash(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	senderSignature := transitionStates[1].InitiatingParticipantSignature
	constraintInput := transitionStates[1].ConstraintInput

	model.SaltedHash = domain.SaltedHash{}

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidInstanceHash(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	constraintInput := transitionStates[1].ConstraintInput

	nextInstance.SaltedHash = domain.SaltedHash{}
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, nextSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidTokenCounts(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[2].Instance
	transition := transitionStates[2].Transition
	senderSignature := transitionStates[2].InitiatingParticipantSignature
	recipientSignature := *transitionStates[2].RespondingParticipantSignature
	constraintInput := transitionStates[2].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidSignature(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	senderSignature := transitionStates[2].InitiatingParticipantSignature
	recipientSignature := *transitionStates[2].RespondingParticipantSignature
	constraintInput := transitionStates[1].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_NotAParticipant(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	constraintInput := transitionStates[1].ConstraintInput
	authentication := circuit.ToAuthentication(nextInstance, transitionStates[1].InitiatingParticipantSignature)
	authentication.MerkleProof.Index = 1

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: authentication,
		RespondingParticipantAuthentication: authentication,
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_AlteredPublicKeys(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	senderSignature := transitionStates[1].InitiatingParticipantSignature
	constraintInput := transitionStates[1].ConstraintInput

	otherPublicKeys := signatureParameters.GetPublicKeys(2)
	currentInstance.PublicKeys[0] = otherPublicKeys[1]
	currentInstance.PublicKeys[1] = otherPublicKeys[0]
	currentInstance.UpdateHash()

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_OverwrittenMessageHash(t *testing.T) {
	model := transitionStates[2].Model
	currentInstance := transitionStates[2].Instance
	nextInstance := transitionStates[3].Instance
	transition := transitionStates[3].Transition
	constraintInput := transitionStates[3].ConstraintInput

	nextInstance.MessageHashes[8] = domain.NewBytesMessage(currentInstance, []byte("Not a purchase order")).Hash.Hash
	nextInstance.UpdateHash()
	senderSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))
	recipientSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(1))

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_OtherParticipant(t *testing.T) {
	model := transitionStates[2].Model
	currentInstance := transitionStates[2].Instance
	nextInstance := transitionStates[3].Instance
	transition := transitionStates[3].Transition
	constraintInput := transitionStates[3].ConstraintInput

	senderSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))
	recipientSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidConstraintInput(t *testing.T) {
	model := transitionStates[2].Model
	currentInstance := transitionStates[2].Instance
	nextInstance := transitionStates[3].Instance
	transition := transitionStates[3].Transition
	senderSignature := transitionStates[3].InitiatingParticipantSignature
	recipientSignature := *transitionStates[3].RespondingParticipantSignature

	constraintInput := transitionStates[3].ConstraintInput
	constraintInput.Messages[0] = domain.NewIntegerMessage(currentInstance, 1)

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidMessageForConstraint(t *testing.T) {
	model := transitionStates[2].Model
	currentInstance := transitionStates[2].Instance
	order := domain.NewIntegerMessage(currentInstance, 6)
	stock := domain.NewIntegerMessage(currentInstance, 4)

	currentInstance.MessageHashes[9] = order.Hash.Hash
	currentInstance.MessageHashes[0] = stock.Hash.Hash
	currentInstance.UpdateHash()
	nextInstance := transitionStates[3].Instance
	nextInstance.MessageHashes[9] = order.Hash.Hash
	nextInstance.MessageHashes[0] = stock.Hash.Hash
	nextInstance.UpdateHash()
	transition := transitionStates[3].Transition
	senderSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(1))
	recipientSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	constraintInput := transitionStates[3].ConstraintInput
	constraintInput.Messages[0] = order

	witness := circuit.TransitionCircuit{
		Model:                               circuit.FromModel(model),
		CurrentInstance:                     circuit.FromInstance(currentInstance),
		NextInstance:                        circuit.FromInstance(nextInstance),
		Transition:                          circuit.ToTransition(model, transition),
		InitiatingParticipantAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		RespondingParticipantAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:                     circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
