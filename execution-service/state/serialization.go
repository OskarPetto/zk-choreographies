package state

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"
)

func (state State) Serialize() domain.Plaintext {
	var modelJson *model.ModelJson = nil
	if state.Model != nil {
		tmp := model.ToJson(*state.Model)
		modelJson = &tmp
	}
	var instanceJson *instance.InstanceJson = nil
	if state.Instance != nil {
		tmp := instance.ToJson(*state.Instance)
		instanceJson = &tmp
	}
	var messageJson *message.MessageJson = nil
	if state.Message != nil {
		tmp := message.ToJson(*state.Message)
		messageJson = &tmp
	}
	stateJson := StateJson{
		Model:    modelJson,
		Instance: instanceJson,
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
	var stateJson StateJson
	err := json.Unmarshal(serializedState.Value, &stateJson)
	if err != nil {
		return State{}, err
	}
	var model *domain.Model = nil
	if stateJson.Model != nil {
		tmp, err := stateJson.Model.ToModel()
		if err != nil {
			return State{}, err
		}
		model = &tmp
	}
	var instance *domain.Instance = nil
	if stateJson.Instance != nil {
		tmp, err := stateJson.Instance.ToInstance()
		if err != nil {
			return State{}, err
		}
		instance = &tmp
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
