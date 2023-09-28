package circuit_test

import (
	"execution-service/circuit"
	"execution-service/domain"
	"execution-service/signature"
	"execution-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var instantiationCircuit circuit.InstantiationCircuit

var signatureService signature.SignatureService = signature.InitializeSignatureService()

func TestInstantiation(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstantiation_InvalidModelHash(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.InvalidHash()),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidInstanceHash(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	instance.Hash = domain.InvalidHash()
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts1(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance2(publicKeys)
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts2(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	instance.TokenCounts[domain.MaxPlaceCount-1] = 0
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidSignature(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	instance2 := testdata.GetModel2Instance2(publicKeys)
	instance2.ComputeHash()
	signature := signatureService.Sign(instance2)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidAuthorization(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 3)
	instance := testdata.GetModel2Instance1([]domain.PublicKey{publicKeys[1], publicKeys[2]})
	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidMessageHashes(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	instance.MessageHashes[0] = domain.HashMessage([]byte("invalid"))

	signature := signatureService.Sign(instance)

	model := testdata.GetModel2()

	witness := circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(domain.HashModel(model)),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
