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

func TestSerializationAndDeserialization(t *testing.T) {
	plainState := state.State{
		Model:    &states[0].Model,
		Instance: &states[0].Instance,
		Message:  states[0].Message,
	}
	result, err := state.Deserialize(plainState.Serialize())
	assert.Nil(t, err)
	assert.Equal(t, result, plainState)
}
