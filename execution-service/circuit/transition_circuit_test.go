package circuit_test

import (
	"proof-service/authentication"
	"proof-service/circuit"
	"proof-service/domain"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit circuit.TransitionCircuit

func TestExecution_NoTokenChange(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
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
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
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
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
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
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance3(publicKeys)
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
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
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
	publicKeys := testdata.GetPublicKeys(3)
	publicKeys = []domain.PublicKey{publicKeys[1], publicKeys[2]}
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
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
	publicKeys := testdata.GetPublicKeys(2)
	publicKeys2 := []domain.PublicKey{publicKeys[1], publicKeys[0]}
	currentInstance := testdata.GetModel2Instance1(publicKeys2)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
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
