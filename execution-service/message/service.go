package message

import (
	"execution-service/domain"
	"execution-service/utils"
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

func (service *MessageService) FindMessageByHashValue(hashValue [domain.HashSize]byte) (domain.Message, error) {
	id := utils.BytesToString(hashValue[:])
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
			message, err := service.FindMessageByHashValue(messageHash)
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
		return fmt.Errorf("message %s has invalid hash", message.String())
	}
	service.messages[message.Hash.String()] = message
	return nil
}
