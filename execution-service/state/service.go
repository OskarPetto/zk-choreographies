package state

import (
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
	serializedState, err := cmd.EncryptedState.Decrypt(privateKey)
	if err != nil {
		return err
	}
	plainState, err := Deserialize(serializedState)
	if err != nil {
		return err
	}
	service.ModelService.ImportModel(plainState.Model)
	service.InstanceService.ImportInstance(plainState.Instance)
	service.MessageService.ImportMessage(*plainState.Message)
	return nil
}
