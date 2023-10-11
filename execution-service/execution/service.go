package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/prover"
	"execution-service/state"
	"fmt"
)

type ExecutionService struct {
	ModelService        model.ModelService
	InstanceService     instance.InstanceService
	MessageService      message.MessageService
	ProverService       prover.IProverService
	SignatureParameters parameters.SignatureParameters
}

func InitializeExecutionService(proverService prover.IProverService) ExecutionService {
	modelService := model.NewModelService()
	instanceService := instance.NewInstanceService()
	messageService := message.NewMessageService()
	signatureParameters := parameters.NewSignatureParameters()
	return NewExecutionService(modelService, instanceService, messageService, proverService, signatureParameters)
}

func NewExecutionService(modelService model.ModelService, instanceService instance.InstanceService, messageService message.MessageService, proverService prover.IProverService, signatureParameters parameters.SignatureParameters) ExecutionService {
	return ExecutionService{
		ModelService:        modelService,
		InstanceService:     instanceService,
		MessageService:      messageService,
		ProverService:       proverService,
		SignatureParameters: signatureParameters,
	}
}

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (ExecutionResult, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ExecutionResult{}, err
	}
	instance, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return ExecutionResult{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveInstantiation(prover.ProveInstantiationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return ExecutionResult{}, err
	}
	service.InstanceService.ImportInstance(instance)
	plainState := domain.NewState(model, instance, nil)
	publicKey := domain.NewPublicKey(privateKey.PublicKey)
	serializedState := state.Serialize(plainState)
	encryptedState := serializedState.Encrypt(privateKey, publicKey)
	return ExecutionResult{
		Instance:       instance,
		Proof:          proof,
		EncryptedState: encryptedState,
	}, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (ExecutionResult, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ExecutionResult{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ExecutionResult{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ExecutionResult{}, err
	}
	if transition.Message != domain.EmptyMessageId && cmd.CreateMessageCommand == nil {
		return ExecutionResult{}, fmt.Errorf("message is required in transition %s", cmd.Transition)
	} else if transition.Message == domain.EmptyMessageId && cmd.CreateMessageCommand != nil {
		return ExecutionResult{}, fmt.Errorf("message is not allowed in transition %s", cmd.Transition)
	}
	nextInstance := currentInstance
	var message *domain.Message = nil
	if transition.Message != domain.EmptyMessageId && cmd.CreateMessageCommand != nil {
		tmp := domain.NewMessage(cmd.CreateMessageCommand.BytesMessage, cmd.CreateMessageCommand.IntegerMessage)
		message = &tmp
		nextInstance = currentInstance.SetMessageHash(transition, message.Hash)
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return ExecutionResult{}, err
	}
	nextInstance, err = nextInstance.ExecuteTransition(transition, constraintInput)
	if err != nil {
		return ExecutionResult{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:           model,
		CurrentInstance: currentInstance,
		NextInstance:    nextInstance,
		Transition:      transition,
		Signature:       nextInstance.Sign(privateKey),
		ConstraintInput: constraintInput,
	})
	if err != nil {
		return ExecutionResult{}, err
	}
	service.InstanceService.ImportInstance(nextInstance)
	plainState := domain.NewState(model, nextInstance, message)
	publicKey := domain.NewPublicKey(privateKey.PublicKey)
	if transition.RespondingParticipant != domain.EmptyParticipantId {
		publicKey = currentInstance.PublicKeys[transition.RespondingParticipant]
	}
	serializedState := state.Serialize(plainState)
	encryptedState := serializedState.Encrypt(privateKey, publicKey)
	return ExecutionResult{
		Instance:       nextInstance,
		Proof:          proof,
		EncryptedState: encryptedState,
	}, nil
}

func (service *ExecutionService) TerminateInstance(cmd TerminateInstanceCommand) (ExecutionResult, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ExecutionResult{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ExecutionResult{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveTermination(prover.ProveTerminationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return ExecutionResult{}, err
	}
	service.InstanceService.DeleteInstance(instance)
	plainState := domain.NewState(model, instance, nil)
	publicKey := domain.NewPublicKey(privateKey.PublicKey)
	serializedState := state.Serialize(plainState)
	encryptedState := serializedState.Encrypt(privateKey, publicKey)
	return ExecutionResult{
		Instance:       instance,
		Proof:          proof,
		EncryptedState: encryptedState,
	}, nil
}
