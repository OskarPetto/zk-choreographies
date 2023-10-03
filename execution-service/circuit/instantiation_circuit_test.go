package circuit_test

import (
	"execution-service/circuit"
	"execution-service/domain"
	"execution-service/parameters"
	"execution-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var instantiationCircuit circuit.InstantiationCircuit

var signatureParameters parameters.SignatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

func TestInstantiation(t *testing.T) {
	model := states[0].Model
	instance := states[0].Instance
	signature := states[0].Signature

	witness := circuit.InstantiationCircuit{
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
	model := states[0].Model
	instance := states[0].Instance
	signature := states[0].Signature

	model.Hash = domain.EmptyHash()

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidInstanceHash(t *testing.T) {
	model := states[0].Model
	instance := states[0].Instance

	instance.Hash = domain.EmptyHash()
	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts1(t *testing.T) {
	model := states[1].Model
	instance := states[1].Instance
	signature := states[1].Signature

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts2(t *testing.T) {
	model := states[len(states)-1].Model
	instance := states[len(states)-1].Instance
	signature := states[len(states)-1].Signature

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidSignature(t *testing.T) {
	model := states[0].Model
	instance := states[0].Instance
	signature := states[1].Signature

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_NotAParticipant(t *testing.T) {
	model := states[0].Model
	instance := states[0].Instance

	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(2))

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidMessageHashes(t *testing.T) {
	model := states[0].Model
	instance := states[0].Instance

	instance.MessageHashes[0] = domain.NewBytesMessage([]byte("invalid")).Hash.Value
	instance.ComputeHash()
	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
