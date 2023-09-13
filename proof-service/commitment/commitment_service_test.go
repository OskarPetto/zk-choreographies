package commitment_test

import (
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
	commitment := commitmentService.CreateCommitment(commitment1.Id, []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	assert.Equal(t, commitment1, commitment)
}
