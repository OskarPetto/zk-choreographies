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

func TestSetMessageHash(t *testing.T) {
	state1 := states[1]
	messageHash := domain.NewBytesMessage([]byte("test")).Hash.Hash
	result := state1.Instance.SetMessageHash(state1.Model.Transitions[2], messageHash)
	assert.Equal(t, state1.Instance.TokenCounts, result.TokenCounts)
	assert.Equal(t, state1.Instance.PublicKeys, result.PublicKeys)
	assert.NotEqual(t, state1.Instance.MessageHashes, result.MessageHashes)
	assert.NotEqual(t, state1.Instance.Hash, result.Hash)
}

func TestExecuteTransition0(t *testing.T) {
	state0 := states[0]
	state1 := states[1]

	constraintInput := domain.EmptyConstraintInput()
	result, err := state0.Instance.ExecuteTransition(state0.Model.Transitions[0], constraintInput)
	assert.Nil(t, err)
	assert.Equal(t, state1.Instance.TokenCounts, result.TokenCounts)
	assert.Equal(t, state1.Instance.PublicKeys, result.PublicKeys)
	assert.Equal(t, state1.Instance.MessageHashes, result.MessageHashes)
	assert.NotEqual(t, state0.Instance.Hash, result.Hash)
}

func TestExecuteTransition2(t *testing.T) {
	state1 := states[1]
	state2 := states[2]
	constraintInput := domain.EmptyConstraintInput()
	result, err := state1.Instance.ExecuteTransition(state1.Model.Transitions[2], constraintInput)
	assert.Nil(t, err)
	assert.Equal(t, state2.Instance.TokenCounts, result.TokenCounts)
	assert.Equal(t, state2.Instance.PublicKeys, result.PublicKeys)
	assert.Equal(t, state1.Instance.MessageHashes, result.MessageHashes)
	assert.NotEqual(t, state1.Instance.Hash, result.Hash)
}
