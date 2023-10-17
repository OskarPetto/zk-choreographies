package message

import (
	"bytes"
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/utils"
	"fmt"
	"sort"
)

type MessageService struct {
	ModelService        model.ModelService
	InstanceService     instance.InstanceService
	messages            map[string]domain.Message
	SignatureParameters parameters.SignatureParameters
}

func NewMessageService(modelService model.ModelService, instanceService instance.InstanceService, signatureParameters parameters.SignatureParameters) MessageService {
	return MessageService{
		messages:            make(map[string]domain.Message),
		ModelService:        modelService,
		InstanceService:     instanceService,
		SignatureParameters: signatureParameters,
	}
}

func (service *MessageService) SaveMessage(message domain.Message) {
	service.messages[message.Id()] = message
}

func (service *MessageService) FindMessagesByInstance(instanceId domain.InstanceId) []domain.Message {
	instance, err := service.InstanceService.FindInstanceById(instanceId)
	if err != nil {
		return []domain.Message{}
	}
	messages := make([]domain.Message, 0, len(service.messages))
	for _, messageHash := range instance.MessageHashes {
		messageId := utils.BytesToString(messageHash.Value[:])
		message, err := service.FindMessageById(messageId)
		if err != nil {
			messages = append(messages, message)
		}
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].CreatedAt > messages[j].CreatedAt
	})
	return messages
}

func (service *MessageService) FindMessageById(id domain.MessageId) (domain.Message, error) {
	message, exists := service.messages[id]
	if !exists {
		return domain.Message{}, fmt.Errorf("message %s not found", id)
	}
	return message, nil
}

func (service *MessageService) FindConstraintInput(constraint domain.Constraint, instance domain.Instance) (domain.ConstraintInput, error) {
	var constraintInput domain.ConstraintInput
	for i, messageId := range constraint.MessageIds {
		coefficient := constraint.Coefficients[i]
		if coefficient != 0 {
			messageHash := instance.MessageHashes[messageId]
			messageId := string(messageHash.Value[:])
			message, err := service.FindMessageById(messageId)
			if err != nil {
				return domain.ConstraintInput{}, err
			}
			constraintInput.Messages[i] = message
		} else {
			constraintInput.Messages[i] = domain.EmptyMessage()
		}
	}
	return constraintInput, nil
}

func (service *MessageService) ImportMessage(cmd ImportMessageCommand) (domain.Signature, error) {
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return domain.Signature{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	ciphertext := cmd.EncryptedMessage
	plaintext, err := ciphertext.Decrypt(privateKey)
	if err != nil {
		return domain.Signature{}, err
	}
	message, err := DeserializeMessage(plaintext)
	if err != nil {
		return domain.Signature{}, err
	}
	if !message.HasValidHash() {
		return domain.Signature{}, fmt.Errorf("message %s has invalid hash", message.Hash.String())
	}
	messageId := domain.EmptyMessageId
	for i, messageHash := range instance.MessageHashes {
		if bytes.Equal(messageHash.Value[:], message.Hash.Hash.Value[:]) {
			messageId = uint16(i)
		}
	}
	if messageId == domain.EmptyMessageId {
		return domain.Signature{}, fmt.Errorf("instance %s does not contain imported message", cmd.Instance)
	}
	service.SaveMessage(message)
	return instance.Sign(privateKey), nil
}

func (service *MessageService) CreateMessage(cmd CreateMessageCommand) (CreateMessageResult, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return CreateMessageResult{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return CreateMessageResult{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return CreateMessageResult{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)

	if transition.Message == domain.EmptyMessageId {
		return CreateMessageResult{}, fmt.Errorf("transition %s has no message", cmd.Transition)
	}
	if transition.Recipient == domain.EmptyParticipantId {
		return CreateMessageResult{}, fmt.Errorf("transition %s has no recipient", cmd.Transition)
	}

	var message domain.Message
	if cmd.BytesMessage != nil {
		message = domain.NewBytesMessage(cmd.BytesMessage)
	} else {
		message = domain.NewIntegerMessage(*cmd.IntegerMessage)
	}
	service.SaveMessage(message)
	nextInstance := instance.SetMessageHash(transition, message.Hash.Hash)
	service.InstanceService.ImportInstance(nextInstance)

	publicKey := instance.FindParticipantById(transition.Recipient)
	plaintext := SerializeMessage(message)
	cipherText := plaintext.Encrypt(privateKey, publicKey)
	return CreateMessageResult{
		Instance:         nextInstance,
		EncryptedMessage: cipherText,
	}, nil
}
