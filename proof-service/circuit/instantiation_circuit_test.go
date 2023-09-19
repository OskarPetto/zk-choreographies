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
	instance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	commitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance1Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   instance,
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
	instance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance1())
	commitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance2Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   instance,
		Commitment: circuit.FromCommitment(commitment),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts1(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance2())
	commitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance2Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   instance,
		Commitment: circuit.FromCommitment(commitment),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts2(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance, _ := circuit.FromInstance(testdata.GetPetriNet1Instance3())
	commitment, _ := commitmentService.CreateCommitment("any", testdata.GetPetriNet1Instance3Serialized())
	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   instance,
		Commitment: circuit.FromCommitment(commitment),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
