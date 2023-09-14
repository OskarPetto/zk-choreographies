package input_test

import (
	"proof-service/circuit/input"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromPetriNet(t *testing.T) {
	petriNet := testdata.GetPetriNet()
	expected := testdata.GetCircuitPetriNet1()
	result := input.FromPetriNet(petriNet)
	assert.Equal(t, expected, result)
}

func TestFromInstance(t *testing.T) {
	instance := testdata.GetInstance1()
	expected := testdata.GetCircuitInstance1()
	result := input.FromInstance(instance)
	assert.Equal(t, expected, result)
}

func TestFromCommitment(t *testing.T) {
	commitment := testdata.GetCommitment1()
	expected := testdata.GetCircuitCommitment1()
	result := input.FromCommitment(commitment)
	assert.Equal(t, expected, result)
}
