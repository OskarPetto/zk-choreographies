package execution_test

import (
	"execution-service/domain"
	"execution-service/execution"
	"execution-service/signature"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ModelPortMock struct {
}

func (service ModelPortMock) FindModelById(id domain.ModelId) (domain.Model, error) {
	return testdata.GetModel2(), nil
}

var signatureService = signature.InitializeSignatureService()
var executionService = execution.InitializeExecutionService()
var instanceService = executionService.InstanceService
var modelService = executionService.ModelService

func TestInstantiateModel(t *testing.T) {
	model := testdata.GetModel2()
	modelService.ImportModel(model)
	instance, err := executionService.InstantiateModel(execution.InstantiateModelCommand{
		Model:      model.Id(),
		PublicKeys: testdata.GetPublicKeys(signatureService, 2),
	})
	assert.Nil(t, err)
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestExecuteTransition(t *testing.T) {
	model := testdata.GetModel2()
	modelService.ImportModel(model)
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	instanceService.ImportInstance(instance)

	result, err := executionService.ExecuteTransition(execution.ExecuteTransitionCommand{
		Model:      model.Id(),
		Instance:   instance.Id(),
		Transition: model.Transitions[0].Id,
	})
	assert.Nil(t, err)
	_, err = instanceService.FindInstanceById(result.Id())
	assert.Nil(t, err)
}
