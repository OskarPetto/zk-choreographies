package circuit_test

import (
	"proof-service/circuit"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var instantiationCircuit circuit.InstantiationCircuit

func TestWithValidWitness(t *testing.T) {
	witness := circuit.InstantiationCircuit{
		Commitment: testdata.GetCircuitCommitment1(),
		Instance:   testdata.GetCircuitInstance1(),
		PetriNet:   testdata.GetCircuitPetriNet1(),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}

}

func TestWithTooManyTokens(t *testing.T) {
	instance := testdata.GetCircuitInstance1()
	instance.TokenCounts[8] = uints.NewU8(1)
	witness := circuit.InstantiationCircuit{
		Commitment: testdata.GetCircuitCommitment1(),
		Instance:   instance,
		PetriNet:   testdata.GetCircuitPetriNet1(),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
