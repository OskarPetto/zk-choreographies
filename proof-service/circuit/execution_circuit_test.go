package circuit_test

import (
	"proof-service/circuit"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var executionCircuit circuit.ExecutionCircuit

func TestNoTokenChange(t *testing.T) {
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   circuit.FromInstance(testdata.GetPetriNet1Instance1()),
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		NextInstance:      circuit.FromInstance(testdata.GetPetriNet1Instance1()),
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		PetriNet:          circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransition0(t *testing.T) {
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   circuit.FromInstance(testdata.GetPetriNet1Instance1()),
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		NextInstance:      circuit.FromInstance(testdata.GetPetriNet1Instance2()),
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		PetriNet:          circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransition7(t *testing.T) {
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   circuit.FromInstance(testdata.GetPetriNet1Instance3()),
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance3Commitment()),
		NextInstance:      circuit.FromInstance(testdata.GetPetriNet1Instance4()),
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance4Commitment()),
		PetriNet:          circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestInvalidCommitments(t *testing.T) {
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   circuit.FromInstance(testdata.GetPetriNet1Instance1()),
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		NextInstance:      circuit.FromInstance(testdata.GetPetriNet1Instance2()),
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		PetriNet:          circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInvalidTransition(t *testing.T) {
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   circuit.FromInstance(testdata.GetPetriNet1Instance2()),
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		NextInstance:      circuit.FromInstance(testdata.GetPetriNet1Instance1()),
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		PetriNet:          circuit.FromPetriNet(testdata.GetPetriNet1()),
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
