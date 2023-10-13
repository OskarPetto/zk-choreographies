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

var instantiationCircuit = circuit.NewInstantiationCircuit()

var signatureParameters parameters.SignatureParameters = parameters.NewSignatureParameters()
var instantiationStates = testdata.GetModel2States(signatureParameters)

func TestInstantiation(t *testing.T) {
	model := instantiationStates[0].Model
	instance := instantiationStates[0].Instance
	signature := instantiationStates[0].Signature

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstantiation_InvalidModelHash(t *testing.T) {
	model := instantiationStates[0].Model
	instance := instantiationStates[0].Instance
	signature := instantiationStates[0].Signature

	model.Hash = domain.SaltedHash{}

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidInstanceHash(t *testing.T) {
	model := instantiationStates[0].Model
	instance := instantiationStates[0].Instance

	instance.Hash = domain.SaltedHash{}
	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts1(t *testing.T) {
	model := instantiationStates[1].Model
	instance := instantiationStates[1].Instance
	signature := instantiationStates[1].Signature

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts2(t *testing.T) {
	model := instantiationStates[len(instantiationStates)-1].Model
	instance := instantiationStates[len(instantiationStates)-1].Instance
	signature := instantiationStates[len(instantiationStates)-1].Signature

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidSignature(t *testing.T) {
	model := instantiationStates[0].Model
	instance := instantiationStates[0].Instance
	signature := instantiationStates[1].Signature

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_NotAParticipant(t *testing.T) {
	model := instantiationStates[0].Model
	instance := instantiationStates[0].Instance

	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(2))

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidMessageHashes(t *testing.T) {
	model := instantiationStates[0].Model
	instance := instantiationStates[0].Instance

	cmd := domain.CreateMessageCommand{BytesMessage: []byte("invalid"), IntegerMessage: nil}
	instance.MessageHashes[0] = domain.CreateMessage(model.Hash.Hash, cmd).Hash.Hash
	instance.UpdateHash()
	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
