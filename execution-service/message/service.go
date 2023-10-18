package message

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/utils"
	"fmt"
	"sort"
)

type MessageService struct {
	InstanceService instance.InstanceService
	messages        map[string]domain.Message
}

func NewMessageService(instanceService instance.InstanceService) MessageService {
	return MessageService{
		messages:        make(map[string]domain.Message),
		InstanceService: instanceService,
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

func (service *MessageService) ImportMessage(message domain.Message) error {
	if !message.HasValidHash() {
		return fmt.Errorf("message %s has invalid hash", message.Hash.String())
	}
	service.SaveMessage(message)
	return nil
}
