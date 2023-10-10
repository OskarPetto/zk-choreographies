package message

import "execution-service/domain"

type SendMessageCommand struct {
	BytesMessage   []byte
	IntegerMessage domain.IntegerType
}
