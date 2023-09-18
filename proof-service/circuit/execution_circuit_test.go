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

func TestExecution_NoTokenChange(t *testing.T) {
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition0(t *testing.T) {
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_InvalidCommitments(t *testing.T) {
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts1(t *testing.T) {
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts2(t *testing.T) {
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance3())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance1Commitment()),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance3Commitment()),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts3(t *testing.T) {
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance4())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.ExecutionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(testdata.GetPetriNet1Instance2Commitment()),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(testdata.GetPetriNet1Instance4Commitment()),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&executionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
