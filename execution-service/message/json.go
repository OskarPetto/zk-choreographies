package message

import (
	"execution-service/domain"
	"execution-service/utils"
)

type SendMessageCommandJson struct {
	IntegerMessage uint   `json:"integerMessage,omitempty"`
	BytesMessage   string `json:"bytesMessage,omitempty"`
}

func (cmd *SendMessageCommandJson) ToMessageCommand() SendMessageCommand {
	bytesMessage, err := utils.StringToBytes(cmd.BytesMessage)
	if err != nil {
		return SendMessageCommand{
			IntegerMessage: int32(cmd.IntegerMessage),
			BytesMessage:   []byte{},
		}
	}

	return SendMessageCommand{
		IntegerMessage: domain.EmptyIntegerMessage,
		BytesMessage:   bytesMessage,
	}
}
