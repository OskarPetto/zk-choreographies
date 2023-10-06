package circuit_test

import (
	"execution-service/circuit"
	"execution-service/domain"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit = circuit.NewTransitionCircuit()

func TestTransition_NoChange(t *testing.T) {
	model := states[4].Model
	currentInstance := states[4].Instance
	transition := states[4].Transition
	signature := states[4].Signature
	constraintInput := states[4].ConstraintInput

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
	for i := 0; i < len(states)-1; i++ {
		model := states[i].Model
		currentInstance := states[i].Instance
		nextInstance := states[i+1].Instance
		transition := states[i+1].Transition
		nextSignature := states[i+1].Signature
		constraintInput := states[i+1].ConstraintInput

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
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	transition := states[1].Transition
	nextSignature := states[1].Signature
	constraintInput := states[1].ConstraintInput

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
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	transition := states[1].Transition
	constraintInput := states[1].ConstraintInput

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
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[2].Instance
	transition := states[2].Transition
	nextSignature := states[2].Signature
	constraintInput := states[2].ConstraintInput

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
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	transition := states[1].Transition
	nextSignature := states[2].Signature
	constraintInput := states[1].ConstraintInput

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
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	transition := states[1].Transition
	constraintInput := states[1].ConstraintInput

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
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	transition := states[1].Transition
	nextSignature := states[1].Signature
	constraintInput := states[1].ConstraintInput

	otherPublicKeys := signatureParameters.GetPublicKeys(3)
	currentInstance.PublicKeys[0] = otherPublicKeys[1]
	currentInstance.PublicKeys[1] = otherPublicKeys[2]
	currentInstance.ComputeHash()

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
	model := states[2].Model
	currentInstance := states[2].Instance
	nextInstance := states[3].Instance
	transition := states[3].Transition
	constraintInput := states[3].ConstraintInput

	nextInstance.MessageHashes[8] = domain.NewBytesMessage([]byte("Not a purchase order")).Hash.Value
	nextInstance.ComputeHash()
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
	model := states[2].Model
	currentInstance := states[2].Instance
	nextInstance := states[3].Instance
	transition := states[3].Transition
	constraintInput := states[3].ConstraintInput

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
	model := states[3].Model
	currentInstance := states[3].Instance
	nextInstance := states[4].Instance
	transition := states[4].Transition
	nextSignature := states[4].Signature

	constraintInput := states[4].ConstraintInput
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
	model := states[3].Model
	currentInstance := states[3].Instance
	message := domain.NewIntegerMessage(6)
	currentInstance.MessageHashes[8] = message.Hash.Value
	currentInstance.ComputeHash()
	nextInstance := states[4].Instance
	nextInstance.MessageHashes[8] = message.Hash.Value
	nextInstance.ComputeHash()
	transition := states[4].Transition
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(1))

	constraintInput := states[4].ConstraintInput
	constraintInput.Messages[0] = message

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
