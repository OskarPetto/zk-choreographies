package instance_test

import (
	"execution-service/instance"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/testdata"
	"execution-service/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters parameters.SignatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

var modelService = model.NewModelService()
var service = instance.NewInstanceService(modelService)

func TestFindInstancesByModel(t *testing.T) {
	instance := states[0].Instance
	service.ImportInstance(instance)
	modelId := utils.BytesToString(instance.Model.Value[:])
	result := service.FindInstancesByModel(modelId)
	assert.Equal(t, 1, len(result))
}

func TestFindInstanceById(t *testing.T) {
	instance := states[0].Instance
	service.ImportInstance(instance)
	_, err := service.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}
