package message

import (
	"bytes"
	"execution-service/domain"
	"execution-service/utils"
	"fmt"
	"sort"
)

type MessageService struct {
	messages map[string]domain.Message
}

func NewMessageService() MessageService {
	return MessageService{
		messages: make(map[string]domain.Message),
	}
}

func (service *MessageService) saveMessage(message domain.Message) {
	service.messages[message.Id()] = message
}

func (service *MessageService) FindMessagesByInstance(instance domain.InstanceId) []domain.Message {
	instanceHash, err := utils.StringToBytes(instance)
	if err != nil {
		return []domain.Message{}
	}
	messages := make([]domain.Message, 0, len(service.messages))
	for _, message := range service.messages {
		if bytes.Equal(message.Instance.Value[:], instanceHash) {
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
	err := message.ValidateHash()
	if err != nil {
		return err
	}
	service.saveMessage(message)
	return nil
}
