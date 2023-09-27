package circuit_test

import (
	"execution-service/authentication"
	"execution-service/domain"
	"execution-service/proof/circuit"
	"execution-service/testdata"
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

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(currentInstance),
		NextInstanceSignature: circuit.FromSignature(signature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
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

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition1(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(nextInstance)

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_InvalidModelHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()
	model.Hash = domain.DefaultHash

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(model),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidInstanceHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
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

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
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

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
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

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
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

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_OverwrittenMessageHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	currentInstance.MessageHashes[8] = domain.HashMessage([]byte("other"))
	currentInstance.ComputeHash()

	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextInstance.ComputeHash()
	nextSignature := signatureService.Sign(nextInstance)

	witness := circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(testdata.GetModel2()),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
