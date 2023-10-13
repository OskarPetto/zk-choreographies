package state_test

import (
	"execution-service/parameters"
	"execution-service/state"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

func TestSerializeAndDeserializeModel(t *testing.T) {
	model := states[2].Model
	result, err := state.DeserializeModel(state.SerializeModel(model))
	assert.Nil(t, err)
	assert.Equal(t, result, model)
}

func TestSerializeAndDeserializeInstance(t *testing.T) {
	instance := states[2].Instance
	result, err := state.DeserializeInstance(state.SerializeInstance(instance))
	assert.Nil(t, err)
	assert.Equal(t, result, instance)
}

func TestSerializeAndDeserializeMessage(t *testing.T) {
	message := *states[2].Message
	result, err := state.DeserializeMessage(state.SerializeMessage(message))
	assert.Nil(t, err)
	assert.Equal(t, result, message)
}
