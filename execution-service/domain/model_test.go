package domain_test

import (
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstantiateModel(t *testing.T) {
	model := testdata.GetModel2()
	publicKeys := testdata.GetPublicKeys(2)
	expected := testdata.GetModel2Instance1(publicKeys)
	instance, err := model.Instantiate(publicKeys)
	assert.Nil(t, err)
	assert.Equal(t, expected.TokenCounts, instance.TokenCounts)
	assert.Equal(t, expected.PublicKeys, instance.PublicKeys)
	assert.Equal(t, expected.MessageHashes, instance.MessageHashes)
	assert.Equal(t, model.Id, instance.Model)
	assert.NotEqual(t, expected.Hash, instance.Hash)
}

func TestFindTransitionById(t *testing.T) {
	id := "ChoreographyTask_0kp4flv_Participant_0x6v44d"
	model := testdata.GetModel2()
	transition, err := model.FindTransitionById(id)
	assert.Nil(t, err)
	assert.Equal(t, id, transition.Id)
}
