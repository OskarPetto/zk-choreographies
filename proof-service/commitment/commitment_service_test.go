package commitment_test

import (
	"fmt"
	"proof-service/commitment"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindCommitment(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	commitment1 := testdata.GetCommitment1()
	commitmentService.SaveCommitment(commitment1)
	commitment2, err := commitmentService.FindCommitment(commitment1.Id)
	assert.Nil(t, err)
	assert.Equal(t, commitment1, commitment2)
}

func TestCreateCommitment(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance := testdata.GetPetriNet1Instance1()
	result := commitmentService.CreateCommitment(instance)
	assert.Equal(t, instance.Id, result.Id)
	fmt.Printf("%+v\n", result)
}
