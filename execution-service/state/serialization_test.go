package state_test

import (
	"execution-service/domain"
	"execution-service/parameters"
	"execution-service/state"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

func TestSerializationAndDeserialization(t *testing.T) {
	plainState := domain.NewState(states[0].Model, states[0].Instance, states[0].Message)
	result, err := state.Deserialize(state.Serialize(plainState))
	assert.Nil(t, err)
	assert.Equal(t, result, plainState)
}
