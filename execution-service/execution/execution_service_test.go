package execution_test

import (
	"proof-service/execution"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstantiatePetriNet(t *testing.T) {
	executionService := execution.NewExecutionService()
	petriNet := testdata.GetPetriNet1()
	publicKeys := testdata.GetPublicKeys(int(petriNet.ParticipantCount))
	instance, err := executionService.InstantiatePetriNet(petriNet, publicKeys)
	assert.Nil(t, err)
	expected := testdata.GetPetriNet1Instance1(publicKeys[0])
	assert.Equal(t, instance.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance.PublicKeys, expected.PublicKeys)
}

func TestExecuteTransition(t *testing.T) {
	executionService := execution.NewExecutionService()
	petriNet := testdata.GetPetriNet1()
	publicKeys := testdata.GetPublicKeys(int(petriNet.ParticipantCount))
	instance, err := executionService.InstantiatePetriNet(petriNet, publicKeys)
	assert.Nil(t, err)
	expected := testdata.GetPetriNet1Instance1(publicKeys[0])
	assert.Equal(t, instance.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance.PublicKeys, expected.PublicKeys)
}
