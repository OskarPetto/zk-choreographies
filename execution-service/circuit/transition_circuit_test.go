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
	signature := transitionStates[4].Signature
	constraintInput := transitionStates[4].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(currentInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(currentInstance, signature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
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
		nextSignature := transitionStates[i+1].Signature
		constraintInput := transitionStates[i+1].ConstraintInput
		witness := circuit.TransitionCircuit{
			Model:           circuit.FromModel(model),
			CurrentInstance: circuit.FromInstance(currentInstance),
			NextInstance:    circuit.FromInstance(nextInstance),
			Transition:      circuit.ToTransition(model, transition),
			Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
			ConstraintInput: circuit.FromConstraintInput(constraintInput),
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
	nextSignature := transitionStates[1].Signature
	constraintInput := transitionStates[1].ConstraintInput

	model.Hash = domain.EmptyHash()

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
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

	nextInstance.Hash = domain.EmptyHash()
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidTokenCounts(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[2].Instance
	transition := transitionStates[2].Transition
	nextSignature := transitionStates[2].Signature
	constraintInput := transitionStates[2].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidSignature(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	nextSignature := transitionStates[2].Signature
	constraintInput := transitionStates[1].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
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

	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(2))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_AlteredPublicKeys(t *testing.T) {
	model := transitionStates[0].Model
	currentInstance := transitionStates[0].Instance
	nextInstance := transitionStates[1].Instance
	transition := transitionStates[1].Transition
	nextSignature := transitionStates[1].Signature
	constraintInput := transitionStates[1].ConstraintInput

	otherPublicKeys := signatureParameters.GetPublicKeys(3)
	currentInstance.PublicKeys[0] = otherPublicKeys[1]
	currentInstance.PublicKeys[1] = otherPublicKeys[2]
	currentInstance.UpdateHash()

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
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

	nextInstance.MessageHashes[8] = domain.NewMessage([]byte("Not a purchase order"), 0).Hash.Value
	nextInstance.UpdateHash()
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
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

	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidConstraintInput(t *testing.T) {
	model := transitionStates[3].Model
	currentInstance := transitionStates[3].Instance
	nextInstance := transitionStates[4].Instance
	transition := transitionStates[4].Transition
	nextSignature := transitionStates[4].Signature

	constraintInput := transitionStates[4].ConstraintInput
	constraintInput.Messages[0].IntegerMessage = 1

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidMessageForConstraint(t *testing.T) {
	model := transitionStates[3].Model
	currentInstance := transitionStates[3].Instance
	order := domain.NewMessage(nil, 6)
	stock := domain.NewMessage(nil, 4)
	currentInstance.MessageHashes[9] = order.Hash.Value
	currentInstance.MessageHashes[0] = stock.Hash.Value
	currentInstance.UpdateHash()
	nextInstance := transitionStates[4].Instance
	nextInstance.MessageHashes[9] = order.Hash.Value
	nextInstance.MessageHashes[0] = stock.Hash.Value
	nextInstance.UpdateHash()
	transition := transitionStates[4].Transition
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(1))

	constraintInput := transitionStates[4].ConstraintInput
	constraintInput.Messages[0] = order

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
