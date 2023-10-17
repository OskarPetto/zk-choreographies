package message_test

import (
	"execution-service/message"
	"execution-service/parameters"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

func TestSerializeAndDeserializeMessage(t *testing.T) {
	domainMessage := *states[2].Message
	result, err := message.DeserializeMessage(message.SerializeMessage(domainMessage))
	assert.Nil(t, err)
	assert.Equal(t, result, domainMessage)
}
