package circuit_test

import (
	"proof-service/circuit"
	"proof-service/commitment"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var transitionCircuit circuit.TransitionCircuit

func TestExecution_NoTokenChange(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	currentCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance1Serialized())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	nextCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance1Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_Transition0(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	currentCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance1Serialized())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	nextCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance2Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecution_InvalidCommitments(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	currentCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance2Serialized())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	nextCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance1Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts1(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	currentCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance2Serialized())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	nextCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance1Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestExecution_InvalidTokenCounts2(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	currentInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	currentCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance1Serialized())
	nextInstance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance3())
	nextCommitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance3Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TransitionCircuit{
		CurrentInstance:   currentInstance,
		CurrentCommitment: circuit.FromCommitment(currentCommitment),
		NextInstance:      nextInstance,
		NextCommitment:    circuit.FromCommitment(nextCommitment),
		PetriNet:          petriNet,
	}

	err := test.IsSolved(&transitionCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
