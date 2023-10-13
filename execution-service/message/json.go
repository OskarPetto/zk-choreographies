package message

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/utils"
)

type MessageJson struct {
	Hash           hash.SaltedHashJson `json:"hash"`
	Model          string              `json:"model"`
	IntegerMessage int                 `json:"integerMessage,omitempty"`
	BytesMessage   string              `json:"bytesMessage,omitempty"`
}

func ToJson(message domain.Message) MessageJson {
	messageJson := MessageJson{
		Hash:  hash.ToJson(message.Hash),
		Model: utils.BytesToString(message.Model.Value[:]),
	}
	if message.IsBytesMessage() {
		messageJson.BytesMessage = utils.BytesToString(message.BytesMessage)
	} else {
		messageJson.IntegerMessage = int(message.IntegerMessage)
	}
	return messageJson
}

func (messageJson *MessageJson) ToMessage() (domain.Message, error) {
	model, err := utils.StringToBytes(messageJson.Model)
	if err != nil {
		return domain.Message{}, err
	}
	hash, err := messageJson.Hash.ToHash()
	if err != nil {
		return domain.EmptyMessage(), err
	}
	message := domain.Message{
		Model: domain.Hash{
			Value: [32]byte(model),
		},
		Hash: hash,
	}
	if messageJson.BytesMessage != "" {
		bytes, err := utils.StringToBytes(messageJson.BytesMessage)
		if err != nil {
			return domain.EmptyMessage(), err
		}
		message.BytesMessage = bytes
	} else {
		message.IntegerMessage = domain.IntegerType(messageJson.IntegerMessage)
	}
	return message, nil
}
