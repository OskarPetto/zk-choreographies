package message

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/instance"
	"execution-service/utils"
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
}

type CreateMessageCommandJson struct {
	Model          string `json:"model"`
	Instance       string `json:"instance"`
	Transition     string `json:"transition"`
	Identity       uint   `json:"identity"`
	IntegerMessage *int   `json:"integerMessage,omitempty"`
	BytesMessage   string `json:"bytesMessage,omitempty"`
}

type CreateMessageResultJson struct {
	Instance         instance.InstanceJson `json:"instance"`
	EncryptedMessage CiphertextJson        `json:"encryptedMessage"`
}

type ImportMessageCommandJson struct {
	EncryptedMessage CiphertextJson `json:"encryptedMessage"`
	Identity         uint           `json:"identity"`
}

func MessageToJson(message domain.Message) MessageJson {
	messageJson := MessageJson{
		Hash: hash.ToJson(message.Hash),
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
		Hash: hash,
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

func (cmd *CreateMessageCommandJson) ToMessageCommand() (CreateMessageCommand, error) {
	createMessageCommand := CreateMessageCommand{
		Model:      cmd.Model,
		Instance:   cmd.Instance,
		Transition: cmd.Transition,
		Identity:   cmd.Identity,
	}
	if cmd.BytesMessage != "" {
		bytes, err := utils.StringToBytes(cmd.BytesMessage)
		if err != nil {
			return CreateMessageCommand{}, err
		}
		createMessageCommand.BytesMessage = bytes
	} else {
		tmp := domain.IntegerType(*cmd.IntegerMessage)
		createMessageCommand.IntegerMessage = &tmp
	}

	return createMessageCommand, nil
}

func (cmd *ImportMessageCommandJson) ToMessageCommand() (ImportMessageCommand, error) {
	ciphertext, err := cmd.EncryptedMessage.toCiphertext()
	if err != nil {
		return ImportMessageCommand{}, err
	}
	return ImportMessageCommand{
		Identity:         cmd.Identity,
		EncryptedMessage: ciphertext,
	}, nil
}

func (json *CiphertextJson) toCiphertext() (domain.Ciphertext, error) {
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

func CreateMessageResultToJson(result CreateMessageResult) CreateMessageResultJson {
	return CreateMessageResultJson{
		Instance:         instance.ToJson(result.Instance),
		EncryptedMessage: toCiphertextJson(result.EncryptedMessage),
	}
}

func toCiphertextJson(encryptedState domain.Ciphertext) CiphertextJson {
	value := utils.BytesToString(encryptedState.Value)
	sender := utils.BytesToString(encryptedState.Sender.Value)
	recipient := utils.BytesToString(encryptedState.Recipient.Value)
	return CiphertextJson{
		Value:     value,
		Sender:    sender,
		Recipient: recipient,
	}
}
