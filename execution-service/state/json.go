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

type encryptedStateJson struct {
	Value  string `json:"value"`
	Sender string `json:"sender"`
}

type ImportStateCommandJson struct {
	EncryptedState encryptedStateJson `json:"encryptedState"`
	Identity       uint               `json:"identity"`
}

func ToJson(encryptedState domain.EncryptedState) encryptedStateJson {
	value := utils.BytesToString(encryptedState.Value)
	sender := utils.BytesToString(encryptedState.Sender.Value)
	return encryptedStateJson{
		Value:  value,
		Sender: sender,
	}
}

func (json *encryptedStateJson) ToEncryptedState() (domain.EncryptedState, error) {
	value, err := utils.StringToBytes(json.Value)
	if err != nil {
		return domain.EncryptedState{}, err
	}
	sender, err := utils.StringToBytes(json.Sender)
	if err != nil {
		return domain.EncryptedState{}, err
	}
	return domain.EncryptedState{
		Value: value,
		Sender: domain.PublicKey{
			Value: sender,
		},
	}, nil
}

func (cmd *ImportStateCommandJson) ToStateCommand() (ImportStateCommand, error) {
	encryptedState, err := cmd.EncryptedState.ToEncryptedState()
	if err != nil {
		return ImportStateCommand{}, err
	}
	return ImportStateCommand{
		EncryptedState: encryptedState,
		Identity:       cmd.Identity,
	}, nil
}
