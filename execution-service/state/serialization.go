package state

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"
)

func Serialize(state State) domain.Plaintext {
	var messageJson *message.MessageJson = nil
	if state.Message != nil {
		tmp := message.ToJson(*state.Message)
		messageJson = &tmp
	}
	stateJson := stateJson{
		Model:    model.ToJson(state.Model),
		Instance: instance.ToJson(state.Instance),
		Message:  messageJson,
	}
	stateJsonBytes, err := json.Marshal(stateJson)
	utils.PanicOnError(err)
	serializedState := domain.Plaintext{
		Value: stateJsonBytes,
	}
	return serializedState
}

func Deserialize(serializedState domain.Plaintext) (State, error) {
	var stateJson stateJson
	err := json.Unmarshal(serializedState.Value, &stateJson)
	if err != nil {
		return State{}, err
	}
	model, err := stateJson.Model.ToModel()
	if err != nil {
		return State{}, err
	}
	instance, err := stateJson.Instance.ToInstance()
	if err != nil {
		return State{}, err
	}
	var message *domain.Message = nil
	if stateJson.Message != nil {
		tmp, err := stateJson.Message.ToMessage()
		if err != nil {
			return State{}, err
		}
		message = &tmp
	}

	return State{
		Model:    model,
		Instance: instance,
		Message:  message,
	}, nil
}
