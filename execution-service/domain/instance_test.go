package domain_test

import (
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecuteTransition0(t *testing.T) {
	model := testdata.GetModel2()
	publicKeys := testdata.GetPublicKeys(2)
	instance1 := testdata.GetModel2Instance1(publicKeys)
	expected := testdata.GetModel2Instance2(publicKeys)
	instance2, err := instance1.ExecuteTransition(model.Transitions[0])
	assert.Nil(t, err)
	assert.Equal(t, instance2.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance2.PublicKeys, expected.PublicKeys)
	assert.Equal(t, instance2.MessageHashes, expected.MessageHashes)
	assert.NotEqual(t, instance2.Hash, instance1.Hash)
}

func TestExecuteTransition1(t *testing.T) {
	model := testdata.GetModel2()
	publicKeys := testdata.GetPublicKeys(2)
	instance1 := testdata.GetModel2Instance2(publicKeys)
	expected := testdata.GetModel2Instance3(publicKeys)
	instance2, err := instance1.ExecuteTransitionWithMessage(model.Transitions[2], []byte("hello"))
	assert.Nil(t, err)
	assert.Equal(t, instance2.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance2.PublicKeys, expected.PublicKeys)
	assert.Equal(t, instance2.MessageHashes, expected.MessageHashes)
	assert.NotEqual(t, instance2.Hash, instance1.Hash)
}
