package domain_test

import (
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstantiateModel(t *testing.T) {
	model := testdata.GetModel2()
	publicKeys := testdata.GetPublicKeys(2)
	expected := testdata.GetModel2Instance1(publicKeys)
	instance, err := model.Instantiate(publicKeys)
	assert.Nil(t, err)
	assert.Equal(t, instance.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance.PublicKeys, expected.PublicKeys)
	assert.NotEqual(t, instance.Hash, expected.Hash)
}

func TestExecuteTransition(t *testing.T) {
	model := testdata.GetModel2()
	publicKeys := testdata.GetPublicKeys(1)
	instance1 := testdata.GetModel2Instance1(publicKeys)
	expected := testdata.GetModel2Instance2(publicKeys)
	instance2, err := instance1.ExecuteTransition(model.Transitions[0])
	assert.Nil(t, err)
	assert.Equal(t, instance2.TokenCounts, expected.TokenCounts)
	assert.Equal(t, instance2.PublicKeys, expected.PublicKeys)
	assert.NotEqual(t, instance2.Hash, instance1.Hash)
}
