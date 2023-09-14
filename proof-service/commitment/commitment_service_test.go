package commitment_test

import (
	"fmt"
	"proof-service/commitment"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var commitment1 = testdata.GetCommitment1()

func TestFindCommitment(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	commitmentService.SaveCommitment(commitment1)
	commitment2, err := commitmentService.FindCommitment(commitment1.Id)
	assert.Nil(t, err)
	assert.Equal(t, commitment1, commitment2)
}

func TestCreateCommitment(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	result, err := commitmentService.CreateCommitment(commitment1.Id, testdata.GetSerializedInstance1())
	assert.Nil(t, err)
	assert.Equal(t, commitment1.Id, result.Id)
	assert.Equal(t, commitment.RandomnessSize, len(result.Randomness))
	fmt.Printf("%+v\n", result)
}
