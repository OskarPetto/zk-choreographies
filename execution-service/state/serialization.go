package state

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"
)

func Serialize(state domain.State) domain.SerializedState {
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
	serializedState := domain.SerializedState{
		Value: stateJsonBytes,
	}
	return serializedState
}

func Deserialize(serializedState domain.SerializedState) (domain.State, error) {
	var stateJson stateJson
	err := json.Unmarshal(serializedState.Value, &stateJson)
	if err != nil {
		return domain.State{}, err
	}
	model, err := stateJson.Model.ToModel()
	if err != nil {
		return domain.State{}, err
	}
	instance, err := stateJson.Instance.ToInstance()
	if err != nil {
		return domain.State{}, err
	}
	var message *domain.Message = nil
	if stateJson.Message != nil {
		tmp, err := stateJson.Message.ToMessage()
		if err != nil {
			return domain.State{}, err
		}
		message = &tmp
	}

	return domain.State{
		Model:    model,
		Instance: instance,
		Message:  message,
	}, nil
}
