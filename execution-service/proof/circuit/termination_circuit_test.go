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

var terminationCircuit circuit.TerminationCircuit

func TestTermination(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestTermination_InvalidModelHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(domain.InvalidHash()),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidInstanceHash(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	instance.Hash = domain.InvalidHash()
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidTokenCounts(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance3(publicKeys)
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidSignature(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	instance2 := testdata.GetModel2Instance2(publicKeys)
	instance2.ComputeHash()
	signature := signatureService.Sign(instance2)

	model := testdata.GetModel2()

	witness := circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidAuthorization(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(3)
	instance := testdata.GetModel2Instance3([]domain.PublicKey{publicKeys[1], publicKeys[2]})
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
