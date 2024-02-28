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
	instanceService := instance.NewInstanceService()
	modelService := model.NewModelService()
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
	if !bytes.Equal(currentInstance.Model.Value[:], model.SaltedHash.Hash.Value[:]) {
		return ExecutedTransitionEvent{}, fmt.Errorf("instance %s is not of model %s", cmd.Instance, cmd.Model)
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	nextInstance, err := currentInstance.ExecuteTransition(transition, constraintInput, nil, nil)
	if err != nil {
		return ExecutedTransitionEvent{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return ExecutedTransitionEvent{}, fmt.Errorf("private key does not exist for identity %d", cmd.Identity)
	}
	senderSignature := nextInstance.Sign(privateKey)
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:                          model,
		CurrentInstance:                currentInstance,
		NextInstance:                   nextInstance,
		Transition:                     transition,
		InitiatingParticipantSignature: senderSignature,
		ConstraintInput:                constraintInput,
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

func (service *ExecutionService) ProveTermination(cmd ProveTerminationCommand) (ProvedTerminationEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ProvedTerminationEvent{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ProvedTerminationEvent{}, err
	}
	if !bytes.Equal(instance.Model.Value[:], model.SaltedHash.Hash.Value[:]) {
		return ProvedTerminationEvent{}, fmt.Errorf("instance %s is not of model %s", cmd.Instance, cmd.Model)
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveTermination(prover.ProveTerminationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return ProvedTerminationEvent{}, err
	}
	return ProvedTerminationEvent{
		Proof: proof,
	}, nil
}

func (service *ExecutionService) CreateInitiatingMessage(cmd CreateInitiatingMessageCommand) (CreatedInitiatingMessageEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return CreatedInitiatingMessageEvent{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return CreatedInitiatingMessageEvent{}, err
	}
	if !bytes.Equal(currentInstance.Model.Value[:], model.SaltedHash.Hash.Value[:]) {
		return CreatedInitiatingMessageEvent{}, fmt.Errorf("instance %s is not of model %s", cmd.Instance, cmd.Model)
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return CreatedInitiatingMessageEvent{}, err
	}
	if transition.InitiatingMessage == domain.EmptyMessageId {
		return CreatedInitiatingMessageEvent{}, fmt.Errorf("transition %s of model %s does not have InitiatingMessage", cmd.Transition, cmd.Model)
	}
	var initiatingMessage domain.Message
	if cmd.BytesMessage != nil {
		initiatingMessage = domain.NewBytesMessage(currentInstance, cmd.BytesMessage)
	} else if cmd.IntegerMessage != nil {
		initiatingMessage = domain.NewIntegerMessage(currentInstance, *cmd.IntegerMessage)
	} else {
		return CreatedInitiatingMessageEvent{}, fmt.Errorf("neither bytes nor integer message was provided for createInitiatingMessage of instance %s", cmd.Instance)
	}
	service.MessageService.ImportMessage(initiatingMessage)
	return CreatedInitiatingMessageEvent{
		Model:              model,
		Instance:           currentInstance,
		Transition:         cmd.Transition,
		InintiatingMessage: initiatingMessage,
	}, nil
}

func (service *ExecutionService) ReceiveInitiatingMessage(cmd ReceiveInitiatingMessageCommand) (ReceivedInitiatingMessageEvent, error) {
	instance := cmd.Instance
	model := cmd.Model
	initiatingMessage := cmd.InitiatingMessage
	if !bytes.Equal(instance.Model.Value[:], model.SaltedHash.Hash.Value[:]) {
		return ReceivedInitiatingMessageEvent{}, fmt.Errorf("instance %s is not of model %s", instance.Id(), model.Id())
	}
	err := service.ModelService.ImportModel(model)
	if err != nil {
		return ReceivedInitiatingMessageEvent{}, err
	}
	err = service.InstanceService.ImportInstance(instance)
	if err != nil {
		return ReceivedInitiatingMessageEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ReceivedInitiatingMessageEvent{}, err
	}
	if !bytes.Equal(initiatingMessage.Instance.Value[:], instance.SaltedHash.Hash.Value[:]) {
		return ReceivedInitiatingMessageEvent{}, fmt.Errorf("message %s is not of instance %s", initiatingMessage.Id(), instance.Id())
	}
	err = service.MessageService.ImportMessage(initiatingMessage)
	if err != nil {
		return ReceivedInitiatingMessageEvent{}, err
	}
	var respondingMessage *domain.Message
	if cmd.BytesMessage != nil {
		tmp := domain.NewBytesMessage(instance, cmd.BytesMessage)
		respondingMessage = &tmp
	} else if cmd.IntegerMessage != nil {
		tmp := domain.NewIntegerMessage(instance, *cmd.IntegerMessage)
		respondingMessage = &tmp
	}
	if respondingMessage != nil {
		service.MessageService.ImportMessage(*respondingMessage)
	}

	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, instance)
	if err != nil {
		return ReceivedInitiatingMessageEvent{}, err
	}
	nextInstance, err := instance.ExecuteTransition(transition, constraintInput, &initiatingMessage, respondingMessage)
	if err != nil {
		return ReceivedInitiatingMessageEvent{}, err
	}

	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return ReceivedInitiatingMessageEvent{}, fmt.Errorf("private key does not exist for identity %d", cmd.Identity)
	}
	respondingParticipantSignature := nextInstance.Sign(privateKey)

	return ReceivedInitiatingMessageEvent{
		Model:                          model.Id(),
		CurrentInstance:                instance.Id(),
		Transition:                     cmd.Transition,
		InitiatingMessage:              initiatingMessage.Id(),
		NextInstance:                   nextInstance,
		RespondingMessage:              respondingMessage,
		RespondingParticipantSignature: respondingParticipantSignature,
	}, nil
}

func (service *ExecutionService) ProveMessageExchange(cmd ProveMessageExchangeCommand) (ProvedMessageExchangeEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ProvedMessageExchangeEvent{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.CurrentInstance)
	if err != nil {
		return ProvedMessageExchangeEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ProvedMessageExchangeEvent{}, err
	}
	nextInstance := cmd.NextInstance
	if !bytes.Equal(nextInstance.Model.Value[:], model.SaltedHash.Hash.Value[:]) {
		return ProvedMessageExchangeEvent{}, fmt.Errorf("next instance %s is not of model %s", cmd.CurrentInstance, cmd.Model)
	}
	err = service.InstanceService.ImportInstance(cmd.NextInstance)
	if err != nil {
		return ProvedMessageExchangeEvent{}, err
	}
	if cmd.RespondingMessage != nil {
		err = service.MessageService.ImportMessage(*cmd.RespondingMessage)
		if err != nil {
			return ProvedMessageExchangeEvent{}, err
		}
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return ProvedMessageExchangeEvent{}, err
	}

	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return ProvedMessageExchangeEvent{}, fmt.Errorf("private key does not exist for identity %d", cmd.Identity)
	}
	initiatingParticipantSignature := nextInstance.Sign(privateKey)
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:                          model,
		CurrentInstance:                currentInstance,
		NextInstance:                   nextInstance,
		Transition:                     transition,
		InitiatingParticipantSignature: initiatingParticipantSignature,
		RespondingParticipantSignature: &cmd.RespondingParticipantSignature,
		ConstraintInput:                constraintInput,
	})
	if err != nil {
		return ProvedMessageExchangeEvent{}, err
	}
	return ProvedMessageExchangeEvent{
		Proof:    proof,
		Instance: nextInstance,
	}, nil

}
