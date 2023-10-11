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
	message.ComputeHash()
	return message
}

func (message *Message) IsBytesMessage() bool {
	return len(message.BytesMessage) > 0
}
