package domain

import "time"

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

func NewBytesMessage(instance Instance, bytes []byte) Message {
	message := Message{
		Instance:     instance.SaltedHash.Hash,
		BytesMessage: bytes,
	}
	message.CreatedAt = time.Now().Unix()
	message.UpdateHash()
	return message
}

func NewIntegerMessage(instance Instance, integer IntegerType) Message {
	message := Message{
		Instance:       instance.SaltedHash.Hash,
		IntegerMessage: integer,
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
