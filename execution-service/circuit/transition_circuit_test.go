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
	signature := transitionStates[4].SenderSignature
	constraintInput := transitionStates[4].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(currentInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(currentInstance, signature),
		RecipientAuthentication: circuit.ToAuthentication(currentInstance, signature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
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
		senderSignature := transitionStates[i+1].SenderSignature
		recipientSignature := transitionStates[i+1].RecipientSignature
		constraintInput := transitionStates[i+1].ConstraintInput

		senderAuthentication := circuit.ToAuthentication(nextInstance, senderSignature)
		recipientAuthentication := senderAuthentication
		if recipientSignature != nil {
			recipientAuthentication = circuit.ToAuthentication(nextInstance, *recipientSignature)
		}
		witness := circuit.TransitionCircuit{
			Model:                   circuit.FromModel(model),
			CurrentInstance:         circuit.FromInstance(currentInstance),
			NextInstance:            circuit.FromInstance(nextInstance),
			Transition:              circuit.ToTransition(model, transition),
			SenderAuthentication:    senderAuthentication,
			RecipientAuthentication: recipientAuthentication,
			ConstraintInput:         circuit.FromConstraintInput(constraintInput),
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
	senderSignature := transitionStates[1].SenderSignature
	constraintInput := transitionStates[1].ConstraintInput

	model.Hash = domain.SaltedHash{}

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
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

	nextInstance.Hash = domain.SaltedHash{}
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, nextSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidTokenCounts(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[2].Instance
	transition := transitionStates[2].Transition
	senderSignature := transitionStates[2].SenderSignature
	recipientSignature := *transitionStates[2].RecipientSignature
	constraintInput := transitionStates[2].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidSignature(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	senderSignature := transitionStates[2].SenderSignature
	recipientSignature := *transitionStates[2].RecipientSignature
	constraintInput := transitionStates[1].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
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
	authentication := circuit.ToAuthentication(nextInstance, transitionStates[1].SenderSignature)
	authentication.MerkleProof.Index = 1

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    authentication,
		RecipientAuthentication: authentication,
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_AlteredPublicKeys(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	senderSignature := transitionStates[1].SenderSignature
	constraintInput := transitionStates[1].ConstraintInput

	otherPublicKeys := signatureParameters.GetPublicKeys(3)
	currentInstance.PublicKeys[0] = otherPublicKeys[1]
	currentInstance.PublicKeys[1] = otherPublicKeys[2]
	currentInstance.UpdateHash()

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, senderSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
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

	nextInstance.MessageHashes[8] = domain.NewBytesMessage([]byte("Not a purchase order")).Hash.Hash
	nextInstance.UpdateHash()
	senderSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))
	recipientSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(1))

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
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
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidConstraintInput(t *testing.T) {
	model := transitionStates[3].Model
	currentInstance := transitionStates[3].Instance
	nextInstance := transitionStates[4].Instance
	transition := transitionStates[4].Transition
	senderSignature := transitionStates[4].SenderSignature
	recipientSignature := *transitionStates[4].RecipientSignature

	constraintInput := transitionStates[4].ConstraintInput
	constraintInput.Messages[0] = domain.NewIntegerMessage(1)

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidMessageForConstraint(t *testing.T) {
	model := transitionStates[3].Model
	currentInstance := transitionStates[3].Instance
	order := domain.NewIntegerMessage(6)
	stock := domain.NewIntegerMessage(4)

	currentInstance.MessageHashes[9] = order.Hash.Hash
	currentInstance.MessageHashes[0] = stock.Hash.Hash
	currentInstance.UpdateHash()
	nextInstance := transitionStates[4].Instance
	nextInstance.MessageHashes[9] = order.Hash.Hash
	nextInstance.MessageHashes[0] = stock.Hash.Hash
	nextInstance.UpdateHash()
	transition := transitionStates[4].Transition
	senderSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(1))
	recipientSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	constraintInput := transitionStates[4].ConstraintInput
	constraintInput.Messages[0] = order

	witness := circuit.TransitionCircuit{
		Model:                   circuit.FromModel(model),
		CurrentInstance:         circuit.FromInstance(currentInstance),
		NextInstance:            circuit.FromInstance(nextInstance),
		Transition:              circuit.ToTransition(model, transition),
		SenderAuthentication:    circuit.ToAuthentication(nextInstance, senderSignature),
		RecipientAuthentication: circuit.ToAuthentication(nextInstance, recipientSignature),
		ConstraintInput:         circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
