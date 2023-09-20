package proof_test

import (
	"proof-service/proof"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProveInstantiation(t *testing.T) {
	proof.LoadParameters()
	instance := testdata.GetPetriNet1Instance1()
	petriNet := testdata.GetPetriNet1()

	proof, err := proof.ProveInstantiation(instance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTransition(t *testing.T) {
	proof.LoadParameters()
	currentInstance := testdata.GetPetriNet1Instance1()
	nextInstance := testdata.GetPetriNet1Instance2()
	petriNet := testdata.GetPetriNet1()

	proof, err := proof.ProveTransition(currentInstance, nextInstance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTermination(t *testing.T) {
	proof.LoadParameters()
	instance := testdata.GetPetriNet1Instance3()
	petriNet := testdata.GetPetriNet1()

	proof, err := proof.ProveTermination(instance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}
