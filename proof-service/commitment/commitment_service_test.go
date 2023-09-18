package commitment_test

import (
	"proof-service/commitment"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var commitment1 = testdata.GetPetriNet1Instance1Commitment()

func TestFindCommitment(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	commitmentService.SaveCommitment(commitment1)
	commitment2, err := commitmentService.FindCommitment(commitment1.Id)
	assert.Nil(t, err)
	assert.Equal(t, commitment1, commitment2)
}

func TestCreateCommitment(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	result, err := commitmentService.CreateCommitment(commitment1.Id, testdata.GetPetriNet1Instance4Serialized())
	assert.Nil(t, err)
	assert.Equal(t, commitment1.Id, result.Id)
	assert.Equal(t, commitment.RandomnessSize, len(result.Randomness))
	//fmt.Printf("%+v\n", result)
}
