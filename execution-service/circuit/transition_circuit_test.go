package circuit_test

import (
	"proof-service/authentication"
	"proof-service/circuit"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit circuit.TransitionCircuit

func TestExecution_NoTokenChange(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetModel2Instance1(publicKey)
	currentInstance.ComputeHash()
	signature := signatureService.Sign(currentInstance)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          currentCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(signature),
		Model:                 model,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition0(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetModel2Instance1(publicKey)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKey)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(nextInstance)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 model,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_InvalidHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetModel2Instance1(publicKey)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKey)
	nextSignature := signatureService.Sign(nextInstance)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 model,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetModel2Instance1(publicKey)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance3(publicKey)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(nextInstance)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 model,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidSignature(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetModel2Instance1(publicKey)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKey)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(currentInstance)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 model,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidAuthorization(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := testdata.GetPublicKeys(2)[1]
	currentInstance := testdata.GetModel2Instance1(publicKey)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKey)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(nextInstance)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 model,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_AlteredPublicKeys(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2)[1])
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKey)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(nextInstance)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 model,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
