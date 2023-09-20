package proof_test

import (
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
	instance := testdata.GetPetriNet1Instance1()
	petriNet := testdata.GetPetriNet1()

	proof, err := proofService.ProveInstantiation(instance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTransition(t *testing.T) {
	currentInstance := testdata.GetPetriNet1Instance1()
	nextInstance := testdata.GetPetriNet1Instance2()
	petriNet := testdata.GetPetriNet1()

	proof, err := proofService.ProveTransition(currentInstance, nextInstance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTermination(t *testing.T) {
	instance := testdata.GetPetriNet1Instance3()
	petriNet := testdata.GetPetriNet1()

	proof, err := proofService.ProveTermination(instance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}
