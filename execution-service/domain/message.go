package domain

import "time"

type MessageId = string

type Message struct {
	Hash           SaltedHash
	IntegerMessage IntegerType
	BytesMessage   []byte
	CreatedAt      int64
}

func EmptyMessage() Message {
	return Message{}
}

func NewBytesMessage(bytes []byte) Message {
	message := Message{
		BytesMessage: bytes,
	}
	message.CreatedAt = time.Now().Unix()
	message.UpdateHash()
	return message
}

func NewIntegerMessage(integer IntegerType) Message {
	message := Message{
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
