package proof_test

import (
	"proof-service/commitment"
	"proof-service/proof"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var proofService proof.ProofService

func TestNewProofService(t *testing.T) {
	proofService = proof.NewProofService()
}

func TestProveInstantiation(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	instance := testdata.GetPetriNet1Instance1()
	commitment := commitmentService.CreateCommitment(instance)
	petriNet := testdata.GetPetriNet1()

	proof, err := proofService.ProveInstantiation(instance, commitment, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTransition(t *testing.T) {
	commitmentService := commitment.NewCommitmentService()
	currentInstance := testdata.GetPetriNet1Instance1()
	currentCommitment := commitmentService.CreateCommitment(currentInstance)
	nextInstance := testdata.GetPetriNet1Instance2()
	nextCommitment := commitmentService.CreateCommitment(nextInstance)
	petriNet := testdata.GetPetriNet1()

	proof, err := proofService.ProveTransition(currentInstance, currentCommitment, nextInstance, nextCommitment, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}
