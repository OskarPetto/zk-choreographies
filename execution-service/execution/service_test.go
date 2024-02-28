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
var states = testdata.GetModel2States(executionService.SignatureParameters)

func TestInitialization(t *testing.T) {
	modelService.ImportModel(states[0].Model)

	for _, state := range states {
		instanceService.ImportInstance(state.Instance)
		if state.InitiatingMessage != nil {
			messageService.ImportMessage(*state.InitiatingMessage)
		}
	}
}

func TestInstantiateModel(t *testing.T) {
	model := states[0].Model
	identity := states[0].InitiatingParticipant
	publicKeys := states[0].Instance.PublicKeys
	event, err := executionService.InstantiateModel(execution.InstantiateModelCommand{
		Model:      model.Id(),
		PublicKeys: publicKeys,
		Identity:   identity,
	})
	assert.Nil(t, err)
	instance := event.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestExecuteTransition0(t *testing.T) {
	model := states[0].Model
	identity := states[0].InitiatingParticipant
	currentInstance := states[0].Instance

	event, err := executionService.ExecuteTransition(execution.ExecuteTransitionCommand{
		Model:      model.Id(),
		Instance:   currentInstance.Id(),
		Transition: model.Transitions[0].Id,
		Identity:   identity,
	})
	assert.Nil(t, err)
	instance := event.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestCreateInitiatingMessageTransition2(t *testing.T) {
	model := states[1].Model
	currentInstance := states[1].Instance

	cmd := execution.CreateInitiatingMessageCommand{
		Model:        model.Id(),
		Instance:     currentInstance.Id(),
		Transition:   model.Transitions[2].Id,
		BytesMessage: []byte("test"),
	}

	event, err := executionService.CreateInitiatingMessage(cmd)
	assert.Nil(t, err)
	assert.Equal(t, cmd.Model, event.Model.Id())
	assert.Equal(t, cmd.Instance, event.Instance.Id())
	assert.Equal(t, cmd.Transition, event.Transition)
	assert.Equal(t, cmd.BytesMessage, event.InintiatingMessage.BytesMessage)
}

func TestReceiveInitiatingMessageTransition2(t *testing.T) {
	currentInstance := states[1].Instance
	model := states[1].Model
	initiatingMessage := *states[2].InitiatingMessage
	integerMessage := states[2].RespondingMessage.IntegerMessage
	cmd := execution.ReceiveInitiatingMessageCommand{
		Model:             model,
		Instance:          currentInstance,
		Transition:        model.Transitions[2].Id,
		Identity:          *states[2].RespondingParticipant,
		InitiatingMessage: initiatingMessage,
		IntegerMessage:    &integerMessage,
	}
	event, err := executionService.ReceiveInitiatingMessage(cmd)
	assert.Nil(t, err)
	assert.Equal(t, cmd.Model.Id(), event.Model)
	assert.Equal(t, cmd.Instance.Id(), event.CurrentInstance)
	assert.Equal(t, cmd.Transition, event.Transition)

	respondingMessage, err := messageService.FindMessageById(event.RespondingMessage.Id())
	assert.Nil(t, err)
	assert.Equal(t, respondingMessage, *event.RespondingMessage)
	assert.Equal(t, integerMessage, event.RespondingMessage.IntegerMessage)
}

func TestProveMessageExchangeTransition2(t *testing.T) {
	currentInstance := states[1].Instance
	model := states[1].Model
	initiatingMessage := *states[2].InitiatingMessage
	cmd := execution.ProveMessageExchangeCommand{
		Model:                          model.Id(),
		CurrentInstance:                currentInstance.Id(),
		Transition:                     model.Transitions[2].Id,
		Identity:                       *states[2].RespondingParticipant,
		InitiatingMessage:              initiatingMessage.Id(),
		NextInstance:                   states[2].Instance,
		RespondingMessage:              states[2].RespondingMessage,
		RespondingParticipantSignature: *states[2].RespondingParticipantSignature,
	}
	event, err := executionService.ProveMessageExchange(cmd)
	assert.Nil(t, err)
	assert.Equal(t, cmd.NextInstance, event.Instance)
}

func TestProveTermination(t *testing.T) {
	model := states[len(states)-1].Model
	instance := states[len(states)-1].Instance
	identity := states[len(states)-1].InitiatingParticipant
	_, err := executionService.ProveTermination(execution.ProveTerminationCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
		Identity: identity,
	})
	assert.Nil(t, err)
}
