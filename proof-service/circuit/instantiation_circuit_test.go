package circuit_test

import (
	"proof-service/circuit"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark/test"
)

var instantiationCircuit circuit.InstantiationCircuit

func TestProverSucceeded(t *testing.T) {
	assert := test.NewAssert(t)

	assert.ProverSucceeded(&instantiationCircuit, &circuit.InstantiationCircuit{
		Commitment: testdata.GetCircuitCommitment1(),
		Instance:   testdata.GetCircuitInstance1(),
		PetriNet:   testdata.GetCircuitPetriNet1(),
	})
}
