package state

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"
)

type CiphertextJson struct {
	Value     string `json:"value"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type StateJson struct {
	Model             *model.ModelJson       `json:"model,omitempty"`
	EncryptedModel    *CiphertextJson        `json:"encryptedModel,omitempty"`
	Instance          *instance.InstanceJson `json:"instance,omitempty"`
	EncryptedInstance *CiphertextJson        `json:"encryptedInstance,omitempty"`
	Message           *message.MessageJson   `json:"message,omitempty"`
	EncryptedMessage  *CiphertextJson        `json:"encryptedMessage,omitempty"`
}

type ImportStateCommandJson struct {
	State    StateJson `json:"state"`
	Identity uint      `json:"identity"`
}

func ToJson(state domain.State) StateJson {
	var modelJson *model.ModelJson = nil
	var encryptedModelJson *CiphertextJson = nil
	var instanceJson *instance.InstanceJson = nil
	var encryptedInstanceJson *CiphertextJson = nil
	var messageJson *message.MessageJson = nil
	var encryptedMessageJson *CiphertextJson = nil

	if state.Model != nil {
		tmp := model.ToJson(*state.Model)
		modelJson = &tmp
	} else if state.EncryptedModel != nil {
		tmp := toCiphertextJson(*state.EncryptedModel)
		encryptedModelJson = &tmp
	}

	if state.Instance != nil {
		tmp := instance.ToJson(*state.Instance)
		instanceJson = &tmp
	} else if state.EncryptedInstance != nil {
		tmp := toCiphertextJson(*state.EncryptedInstance)
		encryptedInstanceJson = &tmp
	}

	if state.Message != nil {
		tmp := message.ToJson(*state.Message)
		messageJson = &tmp
	} else if state.EncryptedMessage != nil {
		tmp := toCiphertextJson(*state.EncryptedMessage)
		encryptedMessageJson = &tmp
	}

	return StateJson{
		Model:             modelJson,
		EncryptedModel:    encryptedModelJson,
		Instance:          instanceJson,
		EncryptedInstance: encryptedInstanceJson,
		Message:           messageJson,
		EncryptedMessage:  encryptedMessageJson,
	}
}

func (stateJson *StateJson) ToState() (domain.State, error) {
	var model *domain.Model = nil
	var encryptedModel *domain.Ciphertext = nil
	var instance *domain.Instance = nil
	var encryptedInstance *domain.Ciphertext = nil
	var message *domain.Message = nil
	var encryptedMessage *domain.Ciphertext = nil

	if stateJson.Model != nil {
		tmp, err := stateJson.Model.ToModel()
		if err != nil {
			return domain.State{}, err
		}
		model = &tmp
	} else if stateJson.EncryptedModel != nil {
		tmp, err := stateJson.EncryptedModel.toCiphertext()
		if err != nil {
			return domain.State{}, err
		}
		encryptedModel = &tmp
	}

	if stateJson.Instance != nil {
		tmp, err := stateJson.Instance.ToInstance()
		if err != nil {
			return domain.State{}, err
		}
		instance = &tmp
	} else if stateJson.EncryptedInstance != nil {
		tmp, err := stateJson.EncryptedInstance.toCiphertext()
		if err != nil {
			return domain.State{}, err
		}
		encryptedInstance = &tmp
	}

	if stateJson.Message != nil {
		tmp, err := stateJson.Message.ToMessage()
		if err != nil {
			return domain.State{}, err
		}
		message = &tmp
	} else if stateJson.EncryptedMessage != nil {
		tmp, err := stateJson.EncryptedMessage.toCiphertext()
		if err != nil {
			return domain.State{}, err
		}
		encryptedMessage = &tmp
	}

	return domain.State{
		Model:             model,
		EncryptedModel:    encryptedModel,
		Instance:          instance,
		EncryptedInstance: encryptedInstance,
		Message:           message,
		EncryptedMessage:  encryptedMessage,
	}, nil
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

func (cmd *ImportStateCommandJson) ToStateCommand() (ImportStateCommand, error) {
	state, err := cmd.State.ToState()
	if err != nil {
		return ImportStateCommand{}, err
	}
	return ImportStateCommand{
		State:    state,
		Identity: cmd.Identity,
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
