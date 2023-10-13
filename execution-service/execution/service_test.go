package execution_test

import (
	"execution-service/domain"
	"execution-service/execution"
	"execution-service/prover"
	"execution-service/state"
	"execution-service/testdata"
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
	state, err := state.Deserialize(*result.PlainState)
	assert.Nil(t, err)
	_, err = instanceService.FindInstanceById(state.Instance.Id())
	assert.Nil(t, err)
	printSize(result)
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
	state, err := state.Deserialize(*result.PlainState)
	assert.Nil(t, err)
	_, err = instanceService.FindInstanceById(state.Instance.Id())
	assert.Nil(t, err)
	printSize(result)
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
			Model:        model.Hash.Hash,
			BytesMessage: []byte("hallo"),
		},
	})
	assert.Nil(t, err)
	plainState, err := state.Deserialize(*result.PlainState)
	assert.Nil(t, err)
	_, err = instanceService.FindInstanceById(plainState.Instance.Id())
	assert.Nil(t, err)

	plaintext, err := result.EncryptedState.Decrypt(privateKey)
	assert.Nil(t, err)
	plainState, err = state.Deserialize(plaintext)
	assert.Nil(t, err)
	_, err = messageService.FindMessageById(plainState.Message.Id())
	assert.Nil(t, err)
	printSize(result)
}

func printSize(result execution.ExecutionResult) {
	plainSize := 0
	encryptedSize := 0
	if result.PlainState != nil {
		plainSize = len(result.PlainState.Value)
	}
	if result.EncryptedState != nil {
		encryptedSize = len(result.EncryptedState.Value)
	}
	fmt.Printf("The length of the plain state is %d bytes\n", plainSize)
	fmt.Printf("The length of the encrypted state is %d bytes\n", encryptedSize)
}
