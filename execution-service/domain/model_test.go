package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstantiateModel(t *testing.T) {
	model := states[0].Model
	publicKeys := signatureParameters.GetPublicKeys(int(model.ParticipantCount))
	state0 := states[0]
	result, err := state0.Model.Instantiate(publicKeys)
	assert.Nil(t, err)
	assert.Equal(t, state0.Instance.TokenCounts, result.TokenCounts)
	assert.Equal(t, state0.Instance.PublicKeys, result.PublicKeys)
	assert.Equal(t, state0.Instance.MessageHashes, result.MessageHashes)
	assert.Equal(t, state0.Model.Hash.Hash, result.Model)
	assert.NotEqual(t, state0.Instance.SaltedHash, result.SaltedHash)
}

func TestFindTransitionById(t *testing.T) {
	id := "ChoreographyTask_0kp4flv"
	model := states[0].Model
	transition, err := model.FindTransitionById(id)
	assert.Nil(t, err)
	assert.Equal(t, id, transition.Id)
}
