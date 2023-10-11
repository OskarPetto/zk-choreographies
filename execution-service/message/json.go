package message

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/utils"
)

type MessageJson struct {
	Hash           hash.HashJson `json:"hash"`
	IntegerMessage int           `json:"integerMessage,omitempty"`
	BytesMessage   string        `json:"bytesMessage,omitempty"`
}

func ToJson(message domain.Message) MessageJson {
	if message.IsBytesMessage() {
		return MessageJson{
			BytesMessage: utils.BytesToString(message.BytesMessage),
			Hash:         hash.ToJson(message.Hash),
		}
	}
	return MessageJson{
		IntegerMessage: int(message.IntegerMessage),
	}
}

func (messageJson *MessageJson) ToMessage() (domain.Message, error) {
	hash, err := messageJson.Hash.ToHash()
	if err != nil {
		return domain.EmptyMessage(), err
	}
	if messageJson.BytesMessage != "" {
		bytes, err := utils.StringToBytes(messageJson.BytesMessage)
		if err != nil {
			return domain.EmptyMessage(), err
		}
		return domain.Message{
			BytesMessage: bytes,
			Hash:         hash,
		}, nil
	}

	return domain.Message{
		IntegerMessage: domain.IntegerType(messageJson.IntegerMessage),
		Hash:           hash,
	}, nil
}
