package state

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
)

type StateService struct {
	ModelService        model.ModelService
	InstanceService     instance.InstanceService
	MessageService      message.MessageService
	SignatureParameters parameters.SignatureParameters
}

func InitializeStateService() StateService {
	modelService := model.NewModelService()
	instanceService := instance.NewInstanceService()
	messageService := message.NewMessageService()
	signatureParameters := parameters.NewSignatureParameters()
	return NewStateService(modelService, instanceService, messageService, signatureParameters)
}

func NewStateService(modelService model.ModelService, instanceService instance.InstanceService, messageService message.MessageService, signatureParameters parameters.SignatureParameters) StateService {
	return StateService{
		ModelService:        modelService,
		InstanceService:     instanceService,
		MessageService:      messageService,
		SignatureParameters: signatureParameters,
	}
}

func (service *StateService) ImportState(cmd ImportStateCommand) error {
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	var model *domain.Model = nil
	var instance *domain.Instance = nil
	var message *domain.Message = nil
	if cmd.State.Model != nil {
		model = cmd.State.Model
	} else if cmd.State.EncryptedModel != nil {
		plainText, err := cmd.State.EncryptedModel.Decrypt(privateKey)
		if err != nil {
			return err
		}
		tmp, err := DeserializeModel(plainText)
		if err != nil {
			return err
		}
		model = &tmp
	}

	if cmd.State.Instance != nil {
		instance = cmd.State.Instance
	} else if cmd.State.EncryptedInstance != nil {
		plainText, err := cmd.State.EncryptedInstance.Decrypt(privateKey)
		if err != nil {
			return err
		}
		tmp, err := DeserializeInstance(plainText)
		if err != nil {
			return err
		}
		instance = &tmp
	}

	if cmd.State.Message != nil {
		message = cmd.State.Message
	} else if cmd.State.EncryptedMessage != nil {
		plainText, err := cmd.State.EncryptedMessage.Decrypt(privateKey)
		if err != nil {
			return err
		}
		tmp, err := DeserializeMessage(plainText)
		if err != nil {
			return err
		}
		message = &tmp
	}

	if model != nil {
		service.ModelService.ImportModel(*model)
	}
	if instance != nil {
		service.InstanceService.ImportInstance(*instance)
	}
	if message != nil {
		service.MessageService.ImportMessage(*message)
	}

	return nil
}
