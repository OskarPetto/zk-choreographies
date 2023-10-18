package message

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/utils"
	"time"
)

type CiphertextJson struct {
	Value     string `json:"value"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type MessageJson struct {
	Hash           hash.SaltedHashJson `json:"hash"`
	IntegerMessage *int                `json:"integerMessage,omitempty"`
	BytesMessage   string              `json:"bytesMessage,omitempty"`
	CreatedAt      time.Time           `json:"createdAt"`
}

type CreateMessageCommandJson struct {
	Identity       uint   `json:"identity"`
	IntegerMessage *int   `json:"integerMessage,omitempty"`
	BytesMessage   string `json:"bytesMessage,omitempty"`
}

type CreateMessageResultJson struct {
	EncryptedMessage CiphertextJson `json:"encryptedMessage"`
}

type ImportMessageCommandJson struct {
	EncryptedMessage CiphertextJson `json:"encryptedMessage"`
	Identity         uint           `json:"identity"`
}

func MessageToJson(message domain.Message) MessageJson {
	messageJson := MessageJson{
		Hash:      hash.ToJson(message.Hash),
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
	message := domain.Message{
		Hash:      hash,
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

func (json *CiphertextJson) ToCiphertext() (domain.Ciphertext, error) {
	value, err := utils.StringToBytes(json.Value)
	if err != nil {
		return domain.Ciphertext{}, err
	}
	sender, err := utils.StringToBytes(json.Sender)
	if err != nil {
		return domain.Ciphertext{}, err
	}
	recipient, err := utils.StringToBytes(json.Recipient)
	if err != nil {
		return domain.Ciphertext{}, err
	}
	return domain.Ciphertext{
		Value: value,
		Sender: domain.PublicKey{
			Value: sender,
		},
		Recipient: domain.PublicKey{
			Value: recipient,
		},
	}, nil
}

func ToCiphertextJson(encryptedState domain.Ciphertext) CiphertextJson {
	value := utils.BytesToString(encryptedState.Value)
	sender := utils.BytesToString(encryptedState.Sender.Value)
	recipient := utils.BytesToString(encryptedState.Recipient.Value)
	return CiphertextJson{
		Value:     value,
		Sender:    sender,
		Recipient: recipient,
	}
}
