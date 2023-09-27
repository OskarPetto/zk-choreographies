package circuit_test

import (
	"proof-service/authentication"
	"proof-service/domain"
	"proof-service/proof/circuit"
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

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(testdata.GetModel2()),
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

	model := testdata.GetModel2()
	model.Hash = domain.DefaultHash

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidInstanceHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	signature := signatureService.Sign(instance)

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(testdata.GetModel2()),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance2(testdata.GetPublicKeys(2))
	instance.ComputeHash()
	signature := signatureService.Sign(instance)

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(testdata.GetModel2()),
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

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(testdata.GetModel2()),
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

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(testdata.GetModel2()),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidMessageHashes(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	instance.MessageHashes[0] = domain.HashMessage([]byte("invalid"))
	instance.ComputeHash()

	signature := signatureService.Sign(instance)

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(testdata.GetModel2()),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
