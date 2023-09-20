package circuit_test

import (
	"proof-service/circuit"
	"proof-service/crypto"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit circuit.TransitionCircuit

func TestExecution_NoTokenChange(t *testing.T) {
	currentInstance := testdata.GetPetriNet1Instance1()
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	currentCommitment := crypto.Commit(currentInstance)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentCircuitInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      currentCircuitInstance,
		NextCommitment:    circuit.FromCommitment(currentCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition0(t *testing.T) {
	currentInstance := testdata.GetPetriNet1Instance1()
	nextInstance := testdata.GetPetriNet1Instance2()
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentCommitment := crypto.Commit(currentInstance)
	nextCommitment := crypto.Commit(nextInstance)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentCircuitInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextCircuitInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_InvalidCommitments(t *testing.T) {
	currentInstance := testdata.GetPetriNet1Instance1()
	nextInstance := testdata.GetPetriNet1Instance2()
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentCommitment := crypto.Commit(nextInstance)
	nextCommitment := crypto.Commit(currentInstance)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentCircuitInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextCircuitInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts1(t *testing.T) {
	currentInstance := testdata.GetPetriNet1Instance1()
	nextInstance := testdata.GetPetriNet1Instance3()
	currentCircuitInstance, _ := circuit.FromInstance(currentInstance)
	nextCircuitInstance, _ := circuit.FromInstance(nextInstance)
	currentCommitment := crypto.Commit(currentInstance)
	nextCommitment := crypto.Commit(nextInstance)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentCircuitInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextCircuitInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
