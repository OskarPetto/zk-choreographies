package circuit_test

import (
	"proof-service/circuit"
	"proof-service/crypto"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit circuit.TransitionCircuit

func TestExecution_NoTokenChange(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	currentSaltedHash := crypto.HashInstance(currentInstance)
	signature := signatureService.Sign(currentSaltedHash)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              currentCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(currentSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(signature),
		PetriNet:                  petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition0(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	nextInstance := testdata.GetPetriNet1Instance2(publicKey)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentSaltedHash := crypto.HashInstance(currentInstance)
	nextSaltedHash := crypto.HashInstance(nextInstance)
	nextSignature := signatureService.Sign(nextSaltedHash)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              nextCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(nextSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(nextSignature),
		PetriNet:                  petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_InvalidSaltedHash(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	nextInstance := testdata.GetPetriNet1Instance2(publicKey)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentSaltedHash := crypto.HashInstance(nextInstance)
	nextSaltedHash := crypto.HashInstance(currentInstance)
	nextSignature := signatureService.Sign(nextSaltedHash)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              nextCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(nextSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(nextSignature),
		PetriNet:                  petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	nextInstance := testdata.GetPetriNet1Instance3(publicKey)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentSaltedHash := crypto.HashInstance(currentInstance)
	nextSaltedHash := crypto.HashInstance(nextInstance)
	nextSignature := signatureService.Sign(nextSaltedHash)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              nextCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(nextSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(nextSignature),
		PetriNet:                  petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidSignature(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	nextInstance := testdata.GetPetriNet1Instance2(publicKey)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentSaltedHash := crypto.HashInstance(currentInstance)
	nextSaltedHash := crypto.HashInstance(nextInstance)
	nextSignature := signatureService.Sign(crypto.HashInstance(nextInstance))
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              nextCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(nextSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(nextSignature),
		PetriNet:                  petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidAuthorization(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := testdata.GetPublicKey1()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	nextInstance := testdata.GetPetriNet1Instance2(publicKey)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentSaltedHash := crypto.HashInstance(currentInstance)
	nextSaltedHash := crypto.HashInstance(nextInstance)
	nextSignature := signatureService.Sign(nextSaltedHash)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              nextCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(nextSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(nextSignature),
		PetriNet:                  petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_AlteredPublicKeys(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(testdata.GetPublicKey1())
	nextInstance := testdata.GetPetriNet1Instance2(publicKey)
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentSaltedHash := crypto.HashInstance(currentInstance)
	nextSaltedHash := crypto.HashInstance(nextInstance)
	nextSignature := signatureService.Sign(nextSaltedHash)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              nextCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(nextSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(nextSignature),
		PetriNet:                  petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
