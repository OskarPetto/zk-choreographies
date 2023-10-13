package execution_test

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/execution"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/prover"
	"execution-service/state"
	"execution-service/testdata"
	"execution-service/utils"
	"fmt"
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
	for _, state := range states {
		modelService.ImportModel(state.Model)
		instanceService.ImportInstance(state.Instance)
		if state.Message != nil {
			messageService.ImportMessage(*state.Message)
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
	instance := result.State.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
	printSize(result.State)
}

func TestExecuteTransition0(t *testing.T) {
	model := states[0].Model
	identity := states[0].Identity
	currentInstance := states[0].Instance

	result, err := executionService.ExecuteTransition(execution.ExecuteTransitionCommand{
		Model:                model.Id(),
		Instance:             currentInstance.Id(),
		Transition:           model.Transitions[0].Id,
		Identity:             identity,
		CreateMessageCommand: nil,
	})
	assert.Nil(t, err)
	instance := result.State.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
	printSize(result.State)
}

func TestExecuteTransition2(t *testing.T) {
	privateKey := executionService.SignatureParameters.GetPrivateKeyForIdentity(1)
	model := states[1].Model
	identity := states[1].Identity
	currentInstance := states[1].Instance

	result, err := executionService.ExecuteTransition(execution.ExecuteTransitionCommand{
		Model:      model.Id(),
		Instance:   currentInstance.Id(),
		Transition: model.Transitions[2].Id,
		Identity:   identity,
		CreateMessageCommand: &domain.CreateMessageCommand{
			BytesMessage: []byte("hallo"),
		},
	})
	assert.Nil(t, err)
	instance := result.State.Instance
	_, err = instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)

	encryptedMessage := result.State.EncryptedMessage
	serializedMessage, err := encryptedMessage.Decrypt(privateKey)
	assert.Nil(t, err)
	message, err := state.DeserializeMessage(serializedMessage)
	assert.Nil(t, err)
	_, err = messageService.FindMessageById(message.Id())
	assert.Nil(t, err)
	printSize(result.State)
}

func TestTerminateInstance(t *testing.T) {
	model := states[len(states)-1].Model
	instance := states[len(states)-1].Instance
	identity := states[len(states)-1].Identity
	result, err := executionService.TerminateInstance(execution.TerminateInstanceCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
		Identity: identity,
	})
	assert.Nil(t, err)
	assert.Equal(t, domain.State{}, result.State)
	printSize(result.State)
}

func printSize(domainState domain.State) {
	plainSize := 0
	encryptedSize := 0
	if domainState.Model != nil {
		value, err := json.Marshal(model.ToJson(*domainState.Model))
		utils.PanicOnError(err)
		plainSize += len(value)
	}
	if domainState.Instance != nil {
		value, err := json.Marshal(instance.ToJson(*domainState.Instance))
		utils.PanicOnError(err)
		plainSize += len(value)
	}
	if domainState.Message != nil {
		value, err := json.Marshal(message.ToJson(*domainState.Message))
		utils.PanicOnError(err)
		plainSize += len(value)
	}
	if domainState.EncryptedModel != nil {
		encryptedSize += len(domainState.EncryptedModel.Value)
		encryptedSize += len(domainState.EncryptedModel.Sender.Value)
		encryptedSize += len(domainState.EncryptedModel.Recipient.Value)
	}
	if domainState.EncryptedInstance != nil {
		encryptedSize += len(domainState.EncryptedInstance.Value)
		encryptedSize += len(domainState.EncryptedInstance.Sender.Value)
		encryptedSize += len(domainState.EncryptedInstance.Recipient.Value)
	}
	if domainState.EncryptedMessage != nil {
		encryptedSize += len(domainState.EncryptedMessage.Value)
		encryptedSize += len(domainState.EncryptedMessage.Sender.Value)
		encryptedSize += len(domainState.EncryptedMessage.Recipient.Value)
	}
	fmt.Printf("The length of the plain state is %d bytes\n", plainSize)
	fmt.Printf("The length of the encrypted state is %d bytes\n", encryptedSize)
}
