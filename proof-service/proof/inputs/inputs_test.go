package inputs_test

import (
	"proof-service/proof/inputs"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromPetriNet(t *testing.T) {
	petriNet := testdata.GetPetriNet1()
	proofPetriNet := testdata.GetProofPetriNet1()
	result, err := inputs.FromPetriNet(petriNet)
	assert.Nil(t, err)
	assert.Equal(t, proofPetriNet, result)
}

func TestFromInstance(t *testing.T) {
	instance := testdata.GetInstance1()
	proofInstance := testdata.GetProofInstance1()
	result, err := inputs.FromInstance(instance)
	assert.Nil(t, err)
	assert.Equal(t, proofInstance, result)
}

func TestFromCommitment(t *testing.T) {
	commitment := testdata.GetCommitment1()
	proofCommitment := testdata.GetProofCommitment1()
	result, err := inputs.FromCommitment(commitment)
	assert.Nil(t, err)
	assert.Equal(t, proofCommitment, result)
}
