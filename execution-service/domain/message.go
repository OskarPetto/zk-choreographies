package domain

type MessageId = string

type Message struct {
	Hash           SaltedHash
	Model          Hash
	IntegerMessage IntegerType
	BytesMessage   []byte
}

type CreateMessageCommand struct {
	BytesMessage   []byte
	IntegerMessage *IntegerType
}

func EmptyMessage() Message {
	return Message{}
}

func CreateMessage(model Hash, cmd CreateMessageCommand) Message {
	message := Message{
		Model: model,
	}
	if len(cmd.BytesMessage) > 0 {
		message.BytesMessage = cmd.BytesMessage
	} else if cmd.IntegerMessage != nil {
		message.IntegerMessage = *cmd.IntegerMessage
	} else {
		return EmptyMessage()
	}
	message.UpdateHash()
	return message
}

func (message Message) Id() MessageId {
	return message.Hash.String()
}

func (message *Message) IsBytesMessage() bool {
	return len(message.BytesMessage) > 0
}
