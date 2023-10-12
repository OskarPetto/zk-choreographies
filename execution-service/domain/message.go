package domain

type Message struct {
	Hash           Hash
	IntegerMessage IntegerType
	BytesMessage   []byte
}

func EmptyMessage() Message {
	return Message{
		Hash: EmptyHash(),
	}
}

func NewMessage(bytesMessage []byte, integerMessage IntegerType) Message {
	var message Message
	if len(bytesMessage) > 0 {
		message = Message{
			BytesMessage: bytesMessage,
		}
	} else {
		message = Message{
			IntegerMessage: integerMessage,
		}
	}
	message.UpdateHash()
	return message
}

func (message *Message) IsBytesMessage() bool {
	return len(message.BytesMessage) > 0
}

func (message Message) String() string {
	if message.IsBytesMessage() {
		return string(message.BytesMessage)
	}
	return string(message.IntegerMessage)
}
