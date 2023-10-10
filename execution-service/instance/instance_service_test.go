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
var publicKeys = signatureParameters.GetPublicKeys(2)

var service = instance.InitializeInstanceService()
var modelService = service.ModelService

func TestFindInstancesByModel(t *testing.T) {
	instance := states[0].Instance
	service.ImportInstance(instance)
	result := service.FindInstancesByModel(instance.Model)
	assert.Equal(t, 1, len(result))
}

func TestFindInstanceById(t *testing.T) {
	instance := states[0].Instance
	service.ImportInstance(instance)
	_, err := service.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestInstantiateModel(t *testing.T) {
	model := states[0].Model
	modelService.ImportModel(model)
	instance, err := service.InstantiateModel(instance.InstantiateModelCommand{
		Model:      model.Id(),
		PublicKeys: publicKeys,
	})
	assert.Nil(t, err)
	_, err = service.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestExecuteTransition2(t *testing.T) {
	model := states[1].Model
	modelService.ImportModel(model)
	currentInstance := states[1].Instance
	service.ImportInstance(currentInstance)

	result, err := service.ExecuteTransition(instance.ExecuteTransitionCommand{
		Model:      model.Id(),
		Instance:   currentInstance.Id(),
		Transition: model.Transitions[2].Id,
	})
	assert.Nil(t, err)
	_, err = service.FindInstanceById(result.Id())
	assert.Nil(t, err)
}
