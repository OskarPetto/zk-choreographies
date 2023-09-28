package circuit_test

import (
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
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	signature := signatureService.Sign(currentInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(currentInstance),
		NextInstanceSignature: circuit.FromSignature(signature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition0(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	nextInstance := testdata.GetModel2Instance2(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition1(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_InvalidModelHash(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.InvalidHash()),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidInstanceHash(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.Hash = domain.InvalidHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidSignature(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	nextInstance := testdata.GetModel2Instance2(publicKeys)
	signature := signatureService.Sign(currentInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(signature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidAuthorization(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 3)
	publicKeys = []domain.PublicKey{publicKeys[1], publicKeys[2]}
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	nextInstance := testdata.GetModel2Instance2(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_AlteredPublicKeys(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	publicKeys2 := []domain.PublicKey{publicKeys[1], publicKeys[0]}
	currentInstance := testdata.GetModel2Instance1(publicKeys2)
	nextInstance := testdata.GetModel2Instance2(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_OverwrittenMessageHash(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	currentInstance.MessageHashes[8] = domain.HashMessage([]byte("other"))

	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)

	model := testdata.GetModel2()

	witness := circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(domain.HashModel(model)),
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
