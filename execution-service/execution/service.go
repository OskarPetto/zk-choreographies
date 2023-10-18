package execution

import (
	"bytes"
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/prover"
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
	instanceService := instance.NewInstanceService(modelService)
	messageService := message.NewMessageService(instanceService)
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

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (InstantiatedModelEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return InstantiatedModelEvent{}, err
	}
	instance, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return InstantiatedModelEvent{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveInstantiation(prover.ProveInstantiationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return InstantiatedModelEvent{}, err
	}
	service.InstanceService.ImportInstance(instance)
	return InstantiatedModelEvent{
		Proof:    proof,
		Instance: instance,
	}, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (ExecutedTransitionEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	nextInstance, err := currentInstance.ExecuteTransition(transition, constraintInput)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	senderSignature := nextInstance.Sign(privateKey)
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:           model,
		CurrentInstance: currentInstance,
		NextInstance:    nextInstance,
		Transition:      transition,
		SenderSignature: senderSignature,
		ConstraintInput: constraintInput,
	})
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	err = service.InstanceService.ImportInstance(nextInstance)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	return ExecutedTransitionEvent{
		Proof:    proof,
		Instance: nextInstance,
	}, nil
}

func (service *ExecutionService) TerminateInstance(cmd TerminateInstanceCommand) (TerminatedInstanceEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return TerminatedInstanceEvent{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return TerminatedInstanceEvent{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveTermination(prover.ProveTerminationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return TerminatedInstanceEvent{}, err
	}
	return TerminatedInstanceEvent{
		Proof: proof,
	}, nil
}

func (service *ExecutionService) SendMessage(cmd SendMessageCommand) (SentMessageEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return SentMessageEvent{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return SentMessageEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return SentMessageEvent{}, err
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return SentMessageEvent{}, err
	}
	if !bytes.Equal(currentInstance.Model.Value[:], model.Hash.Hash.Value[:]) {
		return SentMessageEvent{}, fmt.Errorf("instance %s is not of model %s", cmd.Instance, cmd.Model)
	}
	var messageToSend domain.Message
	if cmd.BytesMessage != nil {
		messageToSend = domain.NewBytesMessage(cmd.BytesMessage)
	} else {
		messageToSend = domain.NewIntegerMessage(*cmd.IntegerMessage)
	}
	nextInstance, err := currentInstance.SendMessage(transition, constraintInput, messageToSend.Hash.Hash)
	if err != nil {
		return SentMessageEvent{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	recipientPublicKey := currentInstance.FindPublicKeyByParticipant(transition.Recipient)
	plaintext := message.SerializeMessage(messageToSend)
	ciphertext := plaintext.Encrypt(privateKey, recipientPublicKey)

	senderSignature := nextInstance.Sign(privateKey)

	service.MessageService.SaveMessage(messageToSend)
	service.InstanceService.SaveInstance(nextInstance)

	return SentMessageEvent{
		Model:            cmd.Model,
		CurrentInstance:  cmd.Instance,
		Transition:       cmd.Transition,
		NextInstance:     nextInstance,
		SenderSignature:  senderSignature,
		EncryptedMessage: ciphertext,
	}, nil
}

func (service *ExecutionService) ReceiveMessage(cmd ReceiveMessageCommand) (ReceivedMessageEvent, error) {
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	plaintext, err := cmd.EncryptedMessage.Decrypt(privateKey)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	receivedMessage, err := message.DeserializeMessage(plaintext)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.CurrentInstance)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	recipientSignature := cmd.NextInstance.Sign(privateKey)
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:              model,
		CurrentInstance:    currentInstance,
		NextInstance:       cmd.NextInstance,
		Transition:         transition,
		SenderSignature:    cmd.SenderSignature,
		RecipientSignature: &recipientSignature,
		ConstraintInput:    constraintInput,
	})
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	err = service.MessageService.ImportMessage(receivedMessage)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	err = service.InstanceService.ImportInstance(cmd.NextInstance)
	if err != nil {
		return ReceivedMessageEvent{}, err
	}
	return ReceivedMessageEvent{
		Proof: proof,
	}, nil
}
