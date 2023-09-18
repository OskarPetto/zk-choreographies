package circuit_test

import (
	"proof-service/circuit"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var instantiationCircuit circuit.InstantiationCircuit

func TestInstantiation(t *testing.T) {
	witness := circuit.InstantiationCircuit{
		Instance:   circuit.FromInstance(testdata.GetPetriNet1Instance1()),
		Commitment: circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		PetriNet:   circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstantiation_InvalidCommitment(t *testing.T) {
	witness := circuit.InstantiationCircuit{
		Instance:   circuit.FromInstance(testdata.GetPetriNet1Instance1()),
		Commitment: circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		PetriNet:   circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts1(t *testing.T) {
	witness := circuit.InstantiationCircuit{
		Instance:   circuit.FromInstance(testdata.GetPetriNet1Instance2()),
		Commitment: circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		PetriNet:   circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts2(t *testing.T) {
	witness := circuit.InstantiationCircuit{
		Instance:   circuit.FromInstance(testdata.GetPetriNet1Instance3()),
		Commitment: circuit.FromCommitment(testdata.GetPetriNet1Instance3Commitment()),
		PetriNet:   circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
