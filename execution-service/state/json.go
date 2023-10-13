package state

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"
)

type StateJson struct {
	Model    *model.ModelJson       `json:"model,omitempty"`
	Instance *instance.InstanceJson `json:"instance,omitempty"`
	Message  *message.MessageJson   `json:"message,omitempty"`
}

type chiphertextJson struct {
	Value     string `json:"value"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type ImportStateCommandJson struct {
	Ciphertext chiphertextJson `json:"encryptedState"`
	Identity   uint            `json:"identity"`
}

func ToCiphertextJson(encryptedState domain.Ciphertext) chiphertextJson {
	value := utils.BytesToString(encryptedState.Value)
	sender := utils.BytesToString(encryptedState.Sender.Value)
	recipient := utils.BytesToString(encryptedState.Recipient.Value)
	return chiphertextJson{
		Value:     value,
		Sender:    sender,
		Recipient: recipient,
	}
}

func (json *chiphertextJson) ToCiphertext() (domain.Ciphertext, error) {
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

func (cmd *ImportStateCommandJson) ToStateCommand() (ImportStateCommand, error) {
	ciphertext, err := cmd.Ciphertext.ToCiphertext()
	if err != nil {
		return ImportStateCommand{}, err
	}
	return ImportStateCommand{
		EncryptedState: ciphertext,
		Identity:       cmd.Identity,
	}, nil
}
