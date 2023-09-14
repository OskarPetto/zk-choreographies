package commitment_test

import (
	"fmt"
	"proof-service/commitment"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var instance = testdata.GetWorkflowInstance1()
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
	commitment, err := commitmentService.CreateCommitment(commitment1.Id, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31})
	assert.Nil(t, err)
	assert.Equal(t, commitment1.Id, commitment.Id)
	assert.Equal(t, 32, len(commitment.Randomness))
	fmt.Printf("%+v\n", commitment)
}
