package execution_test

import (
	"execution-service/execution"
	"execution-service/parameters"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters parameters.SignatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)
var publicKeys = signatureParameters.GetPublicKeys(2)

var executionService = execution.InitializeExecutionService()
var instanceService = executionService.InstanceService
var modelService = executionService.ModelService

func TestInstantiateModel(t *testing.T) {
	model := states[0].Model
	modelService.ImportModel(model)
	instance, err := executionService.InstantiateModel(execution.InstantiateModelCommand{
		Model:      model.Id(),
		PublicKeys: publicKeys,
	})
	assert.Nil(t, err)
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestExecuteTransition2(t *testing.T) {
	model := states[1].Model
	modelService.ImportModel(model)
	instance := states[1].Instance
	instanceService.ImportInstance(instance)

	result, err := executionService.ExecuteTransition(execution.ExecuteTransitionCommand{
		Model:      model.Id(),
		Instance:   instance.Id(),
		Transition: model.Transitions[2].Id,
	})
	assert.Nil(t, err)
	_, err = instanceService.FindInstanceById(result.Id())
	assert.Nil(t, err)
}
