package circuit_test

import (
	"execution-service/circuit"
	"execution-service/domain"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit circuit.TransitionCircuit

func TestTransition_NoChange(t *testing.T) {
	model := states[4].Model
	currentInstance := states[4].Instance
	signature := states[4].Signature
	constraintInput := states[4].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(currentInstance),
		NextSignature:   circuit.FromSignature(signature),
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
		nextSignature := states[i+1].Signature
		constraintInput := states[i+1].ConstraintInput

		witness := circuit.TransitionCircuit{
			Model:           circuit.FromModel(model),
			CurrentInstance: circuit.FromInstance(currentInstance),
			NextInstance:    circuit.FromInstance(nextInstance),
			NextSignature:   circuit.FromSignature(nextSignature),
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
	nextSignature := states[1].Signature
	constraintInput := states[1].ConstraintInput

	model.Hash = domain.EmptyHash()

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidInstanceHash(t *testing.T) {
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	constraintInput := states[1].ConstraintInput

	nextInstance.Hash = domain.EmptyHash()
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidTokenCounts(t *testing.T) {
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[2].Instance
	nextSignature := states[2].Signature
	constraintInput := states[2].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_InvalidSignature(t *testing.T) {
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	nextSignature := states[2].Signature
	constraintInput := states[1].ConstraintInput

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_NotAParticipant(t *testing.T) {
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	constraintInput := states[1].ConstraintInput

	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(2))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_AlteredPublicKeys(t *testing.T) {
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
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
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_OverwrittenMessageHash(t *testing.T) {
	model := states[2].Model
	currentInstance := states[2].Instance
	nextInstance := states[3].Instance
	constraintInput := states[3].ConstraintInput

	nextInstance.MessageHashes[8] = domain.NewBytesMessage([]byte("Not a purchase order")).Hash.Value
	nextInstance.ComputeHash()
	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTransition_OtherParticipant(t *testing.T) {
	model := states[2].Model
	currentInstance := states[2].Instance
	nextInstance := states[3].Instance
	constraintInput := states[3].ConstraintInput

	nextSignature := nextInstance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		NextSignature:   circuit.FromSignature(nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
