package state

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"
)

type stateJson struct {
	Model    model.ModelJson       `json:"model"`
	Instance instance.InstanceJson `json:"instance"`
	Message  *message.MessageJson  `json:"message"`
}

type chiphertextJson struct {
	Value  string `json:"value"`
	Sender string `json:"sender"`
}

type ImportStateCommandJson struct {
	Chiphertext chiphertextJson `json:"encryptedState"`
	Identity    uint            `json:"identity"`
}

func ToJson(encryptedState domain.Chiphertext) chiphertextJson {
	value := utils.BytesToString(encryptedState.Value)
	sender := utils.BytesToString(encryptedState.Sender.Value)
	return chiphertextJson{
		Value:  value,
		Sender: sender,
	}
}

func (json *chiphertextJson) ToChiphertext() (domain.Chiphertext, error) {
	value, err := utils.StringToBytes(json.Value)
	if err != nil {
		return domain.Chiphertext{}, err
	}
	sender, err := utils.StringToBytes(json.Sender)
	if err != nil {
		return domain.Chiphertext{}, err
	}
	return domain.Chiphertext{
		Value: value,
		Sender: domain.PublicKey{
			Value: sender,
		},
	}, nil
}

func (cmd *ImportStateCommandJson) ToStateCommand() (ImportStateCommand, error) {
	ciphertext, err := cmd.Chiphertext.ToChiphertext()
	if err != nil {
		return ImportStateCommand{}, err
	}
	return ImportStateCommand{
		EncryptedState: ciphertext,
		Identity:       cmd.Identity,
	}, nil
}
