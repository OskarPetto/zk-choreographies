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

var instantiationCircuit circuit.InstantiationCircuit

func TestInstantiation(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance := testdata.GetPetriNet1Instance1()
	circuitInstance, _ := circuit.FromInstance(instance)
	commitment := commitmentService.CreateCommitment(instance)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		Commitment: circuit.FromCommitment(commitment),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstantiation_InvalidCommitment(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance := testdata.GetPetriNet1Instance1()
	circuitInstance, _ := circuit.FromInstance(instance)
	commitment := commitmentService.CreateCommitment(testdata.GetPetriNet1Instance2())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		Commitment: circuit.FromCommitment(commitment),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts1(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance := testdata.GetPetriNet1Instance2()
	circuitInstance, _ := circuit.FromInstance(instance)
	commitment := commitmentService.CreateCommitment(instance)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		Commitment: circuit.FromCommitment(commitment),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts2(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance := testdata.GetPetriNet1Instance3()
	circuitInstance, _ := circuit.FromInstance(instance)
	commitment := commitmentService.CreateCommitment(instance)
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		Commitment: circuit.FromCommitment(commitment),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
