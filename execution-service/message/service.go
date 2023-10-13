package message

import (
	"execution-service/domain"
	"fmt"
)

type MessageService struct {
	messages map[string]domain.Message
}

func NewMessageService() MessageService {
	return MessageService{
		messages: make(map[string]domain.Message),
	}
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
	service.messages[message.Id()] = message
	return nil
}
