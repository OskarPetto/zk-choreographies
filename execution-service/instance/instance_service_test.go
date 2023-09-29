package instance_test

import (
	"execution-service/instance"
	"execution-service/parameters"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters parameters.SignatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

func TestFindInstancesByModel(t *testing.T) {
	service := instance.NewInstanceService()
	instance := states[0].Instance
	service.ImportInstance(instance)
	result := service.FindInstancesByModel(instance.Model)
	assert.Equal(t, 1, len(result))
}

func TestFindInstanceById(t *testing.T) {
	service := instance.NewInstanceService()
	instance := states[0].Instance
	service.ImportInstance(instance)
	_, err := service.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}
