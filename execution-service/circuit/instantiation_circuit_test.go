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
	signature := instantiationStates[0].InitiatingParticipantSignature

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
	signature := instantiationStates[0].InitiatingParticipantSignature

	model.SaltedHash = domain.SaltedHash{}

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

	instance.SaltedHash = domain.SaltedHash{}
	pk, _ := signatureParameters.GetPrivateKeyForIdentity(0)
	signature := instance.Sign(pk)

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
	signature := instantiationStates[1].InitiatingParticipantSignature

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
	signature := instantiationStates[len(instantiationStates)-1].InitiatingParticipantSignature

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
	signature := instantiationStates[1].InitiatingParticipantSignature

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
	authentication := circuit.ToAuthentication(instance, instantiationStates[0].InitiatingParticipantSignature)
	authentication.MerkleProof.Index = 1

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: authentication,
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidMessageHashes(t *testing.T) {
	model := instantiationStates[0].Model
	instance := instantiationStates[0].Instance

	instance.MessageHashes[0] = domain.Hash{Value: [32]byte{0, 1, 2, 3}}
	instance.UpdateHash()
	sk, _ := signatureParameters.GetPrivateKeyForIdentity(0)
	signature := instance.Sign(sk)

	witness := circuit.InstantiationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
