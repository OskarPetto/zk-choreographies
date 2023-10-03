package domain

type Message struct {
	Hash           Hash
	IntegerMessage IntegerType
	BytesMessage   []byte
}

func EmptyMessage() Message {
	return Message{}
}

func NewBytesMessage(bytesMessage []byte) Message {
	message := Message{
		BytesMessage: bytesMessage,
	}
	message.ComputeHash()
	return message
}

func NewIntegerMessage(integerMessage IntegerType) Message {
	message := Message{
		IntegerMessage: integerMessage,
	}
	message.ComputeHash()
	return message
}
