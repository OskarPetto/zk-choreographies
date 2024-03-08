package domain_test

import (
	"execution-service/domain"
	"execution-service/parameters"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters parameters.SignatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

func TestExecuteTransition0(t *testing.T) {
	state0 := states[0]
	state1 := states[1]

	conditionInput := domain.EmptyConditionInput()
	result, err := state0.Instance.ExecuteTransition(state0.Model.Transitions[0], conditionInput, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, state1.Instance.TokenCounts, result.TokenCounts)
	assert.Equal(t, state1.Instance.PublicKeys, result.PublicKeys)
	assert.Equal(t, state1.Instance.MessageHashes, result.MessageHashes)
	assert.NotEqual(t, state0.Instance.SaltedHash, result.SaltedHash)
}

func TestExecuteTransition2(t *testing.T) {
	state1 := states[1]
	state2 := states[2]
	transition := state1.Model.Transitions[2]
	initiatingMessage := state2.InitiatingMessage
	respondingMessage := state2.RespondingMessage
	conditionInput := domain.EmptyConditionInput()
	result, err := state1.Instance.ExecuteTransition(transition, conditionInput, initiatingMessage, respondingMessage)
	assert.Nil(t, err)
	assert.Equal(t, state2.Instance.TokenCounts, result.TokenCounts)
	assert.Equal(t, state2.Instance.PublicKeys, result.PublicKeys)
	assert.Equal(t, state2.Instance.MessageHashes, result.MessageHashes)
	assert.NotEqual(t, state1.Instance.SaltedHash, result.SaltedHash)
}

func TestFakeTransition(t *testing.T) {
	state1 := states[1]
	instance1 := state1.Instance
	instance2 := instance1.FakeTransition()
	assert.NotEqual(t, instance1.Id(), instance2.Id())
}
