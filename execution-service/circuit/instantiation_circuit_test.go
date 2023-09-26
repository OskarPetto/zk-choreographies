package circuit_test

import (
	"crypto/sha256"
	"proof-service/authentication"
	"proof-service/circuit"
	"proof-service/domain"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var instantiationCircuit circuit.InstantiationCircuit

func TestInstantiation(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel2())

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

func TestInstantiation_InvalidModelHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel2())
	model.Hash = 1

	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidInstanceHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel2())
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
	instance := testdata.GetModel2Instance2(testdata.GetPublicKeys(2))
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel2())
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
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	instance.ComputeHash()
	instance2 := testdata.GetModel2Instance2(testdata.GetPublicKeys(2))
	instance2.ComputeHash()
	signature := signatureService.Sign(instance2)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel2())
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
	publicKeys := testdata.GetPublicKeys(3)
	instance := testdata.GetModel2Instance1([]domain.PublicKey{publicKeys[1], publicKeys[2]})
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel2())
	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidMessageHashes(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	instance.MessageHashes[0] = domain.MessageHash{
		Value: sha256.Sum256([]byte("invalid")),
	}
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	circuitInstance, _ := circuit.FromInstance(instance)

	model, _ := circuit.FromModel(testdata.GetModel2())

	witness := circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     model,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
