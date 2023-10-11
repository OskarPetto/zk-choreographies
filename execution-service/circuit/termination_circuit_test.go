package circuit_test

import (
	"execution-service/circuit"
	"execution-service/domain"
	"execution-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var terminationCircuit = circuit.NewTerminationCircuit()
var terminationStates = testdata.GetModel2States(signatureParameters)

func TestTermination(t *testing.T) {
	model := terminationStates[len(terminationStates)-1].Model
	instance := terminationStates[len(terminationStates)-1].Instance
	signature := terminationStates[len(terminationStates)-1].Signature

	witness := circuit.TerminationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
		EndPlaceProof:  circuit.ToEndPlaceProof(model, instance),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestTermination_InvalidModelHash(t *testing.T) {
	model := terminationStates[len(terminationStates)-1].Model
	instance := terminationStates[len(terminationStates)-1].Instance
	signature := terminationStates[len(terminationStates)-1].Signature

	model.Hash = domain.EmptyHash()

	witness := circuit.TerminationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
		EndPlaceProof:  circuit.ToEndPlaceProof(model, instance),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidInstanceHash(t *testing.T) {
	model := terminationStates[len(terminationStates)-1].Model
	instance := terminationStates[len(terminationStates)-1].Instance

	instance.Hash = domain.EmptyHash()
	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(0))

	witness := circuit.TerminationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
		EndPlaceProof:  circuit.ToEndPlaceProof(model, instance),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidTokenCounts(t *testing.T) {
	model := terminationStates[len(terminationStates)-2].Model
	instance := terminationStates[len(terminationStates)-2].Instance
	signature := terminationStates[len(terminationStates)-2].Signature

	witness := circuit.TerminationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
		EndPlaceProof:  circuit.ToEndPlaceProof(model, instance),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidSignature(t *testing.T) {
	model := terminationStates[len(terminationStates)-1].Model
	instance := terminationStates[len(terminationStates)-1].Instance
	signature := terminationStates[len(terminationStates)-2].Signature

	witness := circuit.TerminationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
		EndPlaceProof:  circuit.ToEndPlaceProof(model, instance),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_NotAParticipant(t *testing.T) {
	model := terminationStates[len(terminationStates)-1].Model
	instance := terminationStates[len(terminationStates)-1].Instance

	signature := instance.Sign(signatureParameters.GetPrivateKeyForIdentity(2))

	witness := circuit.TerminationCircuit{
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		Model:          circuit.FromModel(model),
		EndPlaceProof:  circuit.ToEndPlaceProof(model, instance),
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
