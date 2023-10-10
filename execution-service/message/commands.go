package message

import "execution-service/domain"

type CreateMessageCommand struct {
	BytesMessage   []byte
	IntegerMessage domain.IntegerType
}
