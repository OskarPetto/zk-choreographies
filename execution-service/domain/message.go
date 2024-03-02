package domain

import (
	"fmt"
	"time"
)

type MessageId = string

type Message struct {
	Hash           SaltedHash
	Instance       Hash
	IntegerMessage IntegerType
	BytesMessage   []byte
	CreatedAt      int64
}

func EmptyMessage() Message {
	return Message{}
}

func NewRespondingBytesMessage(instance Instance, transition Transition, bytes []byte) (Message, error) {
	if transition.RespondingMessage == EmptyMessageId {
		return Message{}, fmt.Errorf("transition %s of model %s does not have an RespondingMessage", transition.Id, instance.Model.String())
	}
	return newBytesMessage(instance, bytes), nil
}

func NewRespondingIntegerMessage(instance Instance, transition Transition, integer IntegerType) (Message, error) {
	if transition.RespondingMessage == EmptyMessageId {
		return Message{}, fmt.Errorf("transition %s of model %s does not have an RespondingMessage", transition.Id, instance.Model.String())
	}
	return newIntegerMessage(instance, integer), nil
}

func NewInitiatingBytesMessage(instance Instance, transition Transition, bytes []byte) (Message, error) {
	if transition.InitiatingMessage == EmptyMessageId {
		return Message{}, fmt.Errorf("transition %s of model %s does not have an InitiatingMessage", transition.Id, instance.Model.String())
	}
	return newBytesMessage(instance, bytes), nil
}

func NewInitiatingIntegerMessage(instance Instance, transition Transition, integer IntegerType) (Message, error) {
	if transition.InitiatingMessage == EmptyMessageId {
		return Message{}, fmt.Errorf("transition %s of model %s does not have an InitiatingMessage", transition.Id, instance.Model.String())
	}
	return newIntegerMessage(instance, integer), nil
}

func newIntegerMessage(instance Instance, integer IntegerType) Message {
	message := Message{
		Instance:       instance.SaltedHash.Hash,
		IntegerMessage: integer,
	}
	message.CreatedAt = time.Now().Unix()
	message.UpdateHash()
	return message
}

func newBytesMessage(instance Instance, bytes []byte) Message {
	message := Message{
		Instance:     instance.SaltedHash.Hash,
		BytesMessage: bytes,
	}
	message.CreatedAt = time.Now().Unix()
	message.UpdateHash()
	return message
}

func (message Message) Id() MessageId {
	return message.Hash.String()
}

func (message *Message) IsBytesMessage() bool {
	return len(message.BytesMessage) > 0
}
