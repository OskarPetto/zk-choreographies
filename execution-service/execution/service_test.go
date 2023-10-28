package execution_test

import (
	"execution-service/execution"
	"execution-service/model"
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
	modelService.ImportModel(model.ImportModelCommand{
		Model:    states[0].Model,
		Instance: states[0].Instance,
	})

	for _, state := range states {
		instanceService.ImportInstance(state.Instance)
		if state.Message != nil {
			messageService.SaveMessage(*state.Message)
		}
	}
}

func TestInstantiateModel(t *testing.T) {
	model := states[0].Model
	identity := states[0].Sender
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
	identity := states[0].Sender
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

func TestSendMessageTransition2(t *testing.T) {
	model := states[1].Model
	identity := states[1].Sender
	currentInstance := states[1].Instance

	cmd := execution.SendMessageCommand{
		Model:        model.Id(),
		Instance:     currentInstance.Id(),
		Transition:   model.Transitions[2].Id,
		Identity:     identity,
		BytesMessage: []byte("test"),
	}

	event, err := executionService.SendMessage(cmd)
	assert.Nil(t, err)
	assert.Equal(t, cmd.Model, event.Model)
	assert.Equal(t, cmd.Instance, event.CurrentInstance)
	assert.Equal(t, cmd.Transition, event.Transition)
	nextInstance := event.NextInstance
	_, err = instanceService.FindInstanceById(nextInstance.Id())
	assert.Nil(t, err)
	transition, err := model.FindTransitionById(cmd.Transition)
	assert.Nil(t, err)
	assert.NotEqual(t, currentInstance.MessageHashes[transition.Message], nextInstance.MessageHashes[transition.Message])
	signature := event.SenderSignature
	assert.True(t, signature.Verify())
	assert.Nil(t, err)
	_, err = messageService.FindMessageById(event.Message.Id())
	assert.Nil(t, err)
}

func TestReceiveMessageTransition2(t *testing.T) {
	currentInstance := states[1].Instance
	model := states[2].Model
	nextInstance := states[2].Instance
	senderSignature := states[2].SenderSignature
	domainMessage := states[2].Message

	cmd := execution.ReceiveMessageCommand{
		Model:           model.Id(),
		CurrentInstance: currentInstance.Id(),
		Transition:      model.Transitions[2].Id,
		Identity:        *states[2].Recipient,
		NextInstance:    nextInstance,
		SenderSignature: senderSignature,
		Message:         domainMessage,
	}
	_, err := executionService.ReceiveMessage(cmd)
	assert.Nil(t, err)
}

func TestTerminateInstance(t *testing.T) {
	model := states[len(states)-1].Model
	instance := states[len(states)-1].Instance
	identity := states[len(states)-1].Sender
	_, err := executionService.TerminateInstance(execution.TerminateInstanceCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
		Identity: identity,
	})
	assert.Nil(t, err)
}
