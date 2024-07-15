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

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (InstanceCreatedEvent, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	instance, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	privateKey, err := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return InstanceCreatedEvent{}, err
	}
	proof, err := service.ProverService.ProveInstantiation(prover.ProveInstantiationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	service.InstanceService.ImportInstance(instance)
	return InstanceCreatedEvent{
		Proof:    proof,
		Instance: instance,
	}, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (InstanceCreatedEvent, error) {
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	model, err := service.ModelService.FindModelById(currentInstance.Model.String())
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	conditionInput, err := service.MessageService.FindConditionInput(transition.Condition, currentInstance)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	nextInstance, err := currentInstance.ExecuteTransition(transition, conditionInput, nil, nil)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	privateKey, err := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return InstanceCreatedEvent{}, err
	}
	senderSignature := nextInstance.Sign(privateKey)
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:                          model,
		CurrentInstance:                currentInstance,
		NextInstance:                   nextInstance,
		Transition:                     transition,
		InitiatingParticipantSignature: senderSignature,
		ConditionInput:                 conditionInput,
	})
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	err = service.InstanceService.ImportInstance(nextInstance)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	return InstanceCreatedEvent{
		Proof:    proof,
		Instance: nextInstance,
	}, nil
}

func (service *ExecutionService) ProveTermination(cmd ProveTerminationCommand) (TerminationProvedEvent, error) {
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return TerminationProvedEvent{}, err
	}
	model, err := service.ModelService.FindModelById(instance.Model.String())
	if err != nil {
		return TerminationProvedEvent{}, err
	}
	privateKey, err := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return TerminationProvedEvent{}, err
	}
	proof, err := service.ProverService.ProveTermination(prover.ProveTerminationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return TerminationProvedEvent{}, err
	}
	return TerminationProvedEvent{
		Proof: proof,
	}, nil
}

func (service *ExecutionService) CreateInitiatingMessage(cmd CreateInitiatingMessageCommand) (InitiatingMessageCreatedEvent, error) {
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return InitiatingMessageCreatedEvent{}, err
	}
	model, err := service.ModelService.FindModelById(currentInstance.Model.String())
	if err != nil {
		return InitiatingMessageCreatedEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return InitiatingMessageCreatedEvent{}, err
	}
	var initiatingMessage *domain.Message = nil
	if cmd.BytesMessage != nil {
		tmp, err := domain.NewInitiatingBytesMessage(currentInstance, transition, cmd.BytesMessage)
		if err != nil {
			return InitiatingMessageCreatedEvent{}, err
		}
		initiatingMessage = &tmp
	} else if cmd.IntegerMessage != nil {
		tmp, err := domain.NewInitiatingIntegerMessage(currentInstance, transition, *cmd.IntegerMessage)
		if err != nil {
			return InitiatingMessageCreatedEvent{}, err
		}
		initiatingMessage = &tmp
	}
	if initiatingMessage != nil {
		service.MessageService.ImportMessage(*initiatingMessage)
	}
	return InitiatingMessageCreatedEvent{
		Model:              model,
		Instance:           currentInstance,
		Transition:         cmd.Transition,
		InintiatingMessage: initiatingMessage,
	}, nil
}

func (service *ExecutionService) ReceiveInitiatingMessage(cmd ReceiveInitiatingMessageCommand) (InitiatingMessageReceivedEvent, error) {
	instance := cmd.Instance
	model := cmd.Model
	initiatingMessage := cmd.InitiatingMessage
	if !bytes.Equal(instance.Model.Value[:], model.SaltedHash.Hash.Value[:]) {
		return InitiatingMessageReceivedEvent{}, fmt.Errorf("instance %s is not of model %s", instance.Id(), model.Id())
	}
	if initiatingMessage != nil && !bytes.Equal(initiatingMessage.Instance.Value[:], instance.SaltedHash.Hash.Value[:]) {
		return InitiatingMessageReceivedEvent{}, fmt.Errorf("message %s is not of instance %s", initiatingMessage.Id(), instance.Id())
	}
	err := service.ModelService.ImportModel(model)
	if err != nil {
		return InitiatingMessageReceivedEvent{}, err
	}
	err = service.InstanceService.ImportInstance(instance)
	if err != nil {
		return InitiatingMessageReceivedEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return InitiatingMessageReceivedEvent{}, err
	}
	var initiatingMessageId *string = nil
	if initiatingMessage != nil {
		err = service.MessageService.ImportMessage(*initiatingMessage)
		if err != nil {
			return InitiatingMessageReceivedEvent{}, err
		}
		tmp := initiatingMessage.Id()
		initiatingMessageId = &tmp
	}
	var respondingMessage *domain.Message
	if cmd.BytesMessage != nil {
		tmp, err := domain.NewInitiatingBytesMessage(instance, transition, cmd.BytesMessage)
		if err != nil {
			return InitiatingMessageReceivedEvent{}, err
		}
		respondingMessage = &tmp
	} else if cmd.IntegerMessage != nil {
		tmp, err := domain.NewInitiatingIntegerMessage(instance, transition, *cmd.IntegerMessage)
		if err != nil {
			return InitiatingMessageReceivedEvent{}, err
		}
		respondingMessage = &tmp
	}
	if respondingMessage != nil {
		service.MessageService.ImportMessage(*respondingMessage)
	}

	conditionInput, err := service.MessageService.FindConditionInput(transition.Condition, instance)
	if err != nil {
		return InitiatingMessageReceivedEvent{}, err
	}
	nextInstance, err := instance.ExecuteTransition(transition, conditionInput, initiatingMessage, respondingMessage)
	if err != nil {
		return InitiatingMessageReceivedEvent{}, err
	}

	privateKey, err := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return InitiatingMessageReceivedEvent{}, err
	}
	respondingParticipantSignature := nextInstance.Sign(privateKey)

	if initiatingMessage != nil {

	}

	return InitiatingMessageReceivedEvent{
		Model:                          model.Id(),
		CurrentInstance:                instance.Id(),
		Transition:                     cmd.Transition,
		InitiatingMessage:              initiatingMessageId,
		NextInstance:                   nextInstance,
		RespondingMessage:              respondingMessage,
		RespondingParticipantSignature: respondingParticipantSignature,
	}, nil
}

func (service *ExecutionService) ProveMessageExchange(cmd ProveMessageExchangeCommand) (InstanceCreatedEvent, error) {
	nextInstance := cmd.NextInstance
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.CurrentInstance)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	model, err := service.ModelService.FindModelById(currentInstance.Model.String())
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	var initiatingMessage *domain.Message = nil
	if cmd.InitiatingMessage != nil {
		tmp, err := service.MessageService.FindMessageById(*cmd.InitiatingMessage)
		if err != nil {
			return InstanceCreatedEvent{}, err
		}
		initiatingMessage = &tmp
	}

	err = cmd.NextInstance.ValidateMessages(transition, initiatingMessage, cmd.RespondingMessage)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	err = service.InstanceService.ImportInstance(cmd.NextInstance)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	if cmd.RespondingMessage != nil {
		err = service.MessageService.ImportMessage(*cmd.RespondingMessage)
		if err != nil {
			return InstanceCreatedEvent{}, err
		}
	}
	conditionInput, err := service.MessageService.FindConditionInput(transition.Condition, currentInstance)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	privateKey, err := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return InstanceCreatedEvent{}, err
	}
	initiatingParticipantSignature := nextInstance.Sign(privateKey)
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:                          model,
		CurrentInstance:                currentInstance,
		NextInstance:                   nextInstance,
		Transition:                     transition,
		InitiatingParticipantSignature: initiatingParticipantSignature,
		RespondingParticipantSignature: &cmd.RespondingParticipantSignature,
		ConditionInput:                 conditionInput,
	})
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	return InstanceCreatedEvent{
		Proof:    proof,
		Instance: nextInstance,
	}, nil

}

func (service *ExecutionService) FakeTransition(cmd FakeTransitionCommand) (InstanceCreatedEvent, error) {
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	model, err := service.ModelService.FindModelById(instance.Model.String())
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	privateKey, err := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	if privateKey == nil {
		return InstanceCreatedEvent{}, err
	}
	instanceWithDifferentHash := instance.FakeTransition()
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:                          model,
		CurrentInstance:                instance,
		NextInstance:                   instanceWithDifferentHash,
		Transition:                     model.Transitions[0],
		InitiatingParticipantSignature: instanceWithDifferentHash.Sign(privateKey),
		ConditionInput:                 domain.EmptyConditionInput(),
	})
	if err != nil {
		return InstanceCreatedEvent{}, err
	}
	return InstanceCreatedEvent{
		Instance: instanceWithDifferentHash,
		Proof:    proof,
	}, nil
}
