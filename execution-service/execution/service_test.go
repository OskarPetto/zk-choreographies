package execution_test

import (
	"execution-service/execution"
	"execution-service/prover"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockProverService struct {
}

func (service MockProverService) ProveInstantiation(cmd prover.ProveInstantiationCommand) (prover.Proof, error) {
	return prover.Proof{}, nil
}

func (service MockProverService) ProveTransition(cmd prover.ProveTransitionCommand) (prover.Proof, error) {
	return prover.Proof{}, nil

}
func (service MockProverService) ProveTermination(cmd prover.ProveTerminationCommand) (prover.Proof, error) {
	return prover.Proof{}, nil
}

var executionService = execution.InitializeExecutionService(MockProverService{})
var modelService = executionService.ModelService
var instanceService = executionService.InstanceService
var messageService = executionService.MessageService
var signatureService = executionService.SignatureService
var states = testdata.GetModel2States(executionService.SignatureParameters)

func TestInitialization(t *testing.T) {
	for _, state := range states {
		modelService.ImportModel(state.Model)
		instanceService.ImportInstance(state.Instance)
		if state.Message != nil {
			messageService.SaveMessage(*state.Message)
		}
		if state.RecipientSignature != nil {
			signatureService.ImportSignature(*state.RecipientSignature)
		}
	}
}

func TestInstantiateModel(t *testing.T) {
	model := states[0].Model
	identity := states[0].Identity
	publicKeys := states[0].Instance.PublicKeys
	result, err := executionService.InstantiateModel(execution.InstantiateModelCommand{
		Model:      model.Id(),
		PublicKeys: publicKeys,
		Identity:   identity,
	})
	assert.Nil(t, err)
	instance := result.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestExecuteTransition0(t *testing.T) {
	model := states[0].Model
	identity := states[0].Identity
	currentInstance := states[0].Instance

	result, err := executionService.ExecuteTransition(execution.ExecuteTransitionCommand{
		Model:      model.Id(),
		Instance:   currentInstance.Id(),
		Transition: model.Transitions[0].Id,
		Identity:   identity,
	})
	assert.Nil(t, err)
	instance := result.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestExecuteTransition2(t *testing.T) {
	model := states[1].Model
	identity := states[1].Identity
	currentInstance := states[1].Instance

	result, err := executionService.ExecuteTransition(execution.ExecuteTransitionCommand{
		Model:      model.Id(),
		Instance:   currentInstance.Id(),
		Transition: model.Transitions[2].Id,
		Identity:   identity,
	})
	assert.Nil(t, err)
	instance := result.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestTerminateInstance(t *testing.T) {
	model := states[len(states)-1].Model
	instance := states[len(states)-1].Instance
	identity := states[len(states)-1].Identity
	_, err := executionService.TerminateInstance(execution.TerminateInstanceCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
		Identity: identity,
	})
	assert.Nil(t, err)
}
