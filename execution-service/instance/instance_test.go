package instance_test

import (
	"proof-service/instance"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstantiatePetriNet(t *testing.T) {
	petriNet := testdata.GetPetriNet1()
	publicKeys := testdata.GetPublicKeys(1)
	expected := testdata.GetPetriNet1Instance1(publicKeys[0])
	instance, err := instance.InstantiatePetriNet(petriNet, publicKeys)
	assert.Nil(t, err)
	assert.Equal(t, instance.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance.PublicKeys, expected.PublicKeys)
	assert.NotEqual(t, instance.Hash, expected.Hash)
}

func TestExecuteTransition(t *testing.T) {
	petriNet := testdata.GetPetriNet1()
	publicKey := testdata.GetPublicKeys(1)[0]
	instance1 := testdata.GetPetriNet1Instance1(publicKey)
	expected := testdata.GetPetriNet1Instance2(publicKey)
	instance2, err := instance1.ExecuteTransition(petriNet.Transitions[0])
	assert.Nil(t, err)
	assert.Equal(t, instance2.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance2.PublicKeys, expected.PublicKeys)
	assert.NotEqual(t, instance2.Hash, instance1.Hash)
}
