package message

import (
	"execution-service/domain"
	"execution-service/utils"
)

type CreateMessageCommandJson struct {
	IntegerMessage uint   `json:"integerMessage,omitempty"`
	BytesMessage   string `json:"bytesMessage,omitempty"`
}

func (cmd *CreateMessageCommandJson) ToMessageCommand() CreateMessageCommand {
	bytesMessage, err := utils.StringToBytes(cmd.BytesMessage)
	if err != nil {
		return CreateMessageCommand{
			IntegerMessage: int32(cmd.IntegerMessage),
			BytesMessage:   []byte{},
		}
	}

	return CreateMessageCommand{
		IntegerMessage: domain.EmptyIntegerMessage,
		BytesMessage:   bytesMessage,
	}
}
