package message

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/utils"
	"time"
)

type MessageJson struct {
	Hash           hash.SaltedHashJson `json:"hash"`
	Instance       string              `json:"instance"`
	IntegerMessage *int                `json:"integerMessage,omitempty"`
	BytesMessage   string              `json:"bytesMessage,omitempty"`
	CreatedAt      time.Time           `json:"createdAt"`
}

func ToJson(message domain.Message) MessageJson {
	messageJson := MessageJson{
		Hash:      hash.ToJson(message.Hash),
		Instance:  utils.BytesToString(message.Instance.Value[:]),
		CreatedAt: time.Unix(message.CreatedAt, 0),
	}
	if message.IsBytesMessage() {
		messageJson.BytesMessage = utils.BytesToString(message.BytesMessage)
	} else {
		tmp := int(message.IntegerMessage)
		messageJson.IntegerMessage = &tmp
	}
	return messageJson
}

func (messageJson *MessageJson) ToMessage() (domain.Message, error) {
	hash, err := messageJson.Hash.ToHash()
	if err != nil {
		return domain.EmptyMessage(), err
	}
	instance, err := utils.StringToBytes(messageJson.Instance)
	if err != nil {
		return domain.EmptyMessage(), err
	}
	instanceFixedSize := [domain.HashSize]byte(instance)
	message := domain.Message{
		Hash: hash,
		Instance: domain.Hash{
			Value: instanceFixedSize,
		},
		CreatedAt: messageJson.CreatedAt.Unix(),
	}
	if messageJson.BytesMessage != "" {
		bytes, err := utils.StringToBytes(messageJson.BytesMessage)
		if err != nil {
			return domain.EmptyMessage(), err
		}
		message.BytesMessage = bytes
	} else {
		message.IntegerMessage = domain.IntegerType(*messageJson.IntegerMessage)
	}
	return message, nil
}
