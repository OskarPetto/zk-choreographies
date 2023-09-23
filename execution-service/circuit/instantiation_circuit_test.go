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

var instantiationCircuit circuit.InstantiationCircuit

func TestInstantiation(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetModel1Instance1(publicKey)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel1())
	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstantiation_InvalidHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetModel1Instance1(publicKey)
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel1())
	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetModel1Instance2(publicKey)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel1())
	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidSignature(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetModel1Instance1(publicKey)
	instance.ComputeHash()
	instance2 := testdata.GetModel1Instance2(publicKey)
	instance2.ComputeHash()
	signature := signatureService.Sign(instance2)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel1())
	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidAuthorization(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := testdata.GetPublicKeys(2)[1]
	instance := testdata.GetModel1Instance1(publicKey)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel1())
	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
