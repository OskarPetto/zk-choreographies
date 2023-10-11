package state

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

type State struct {
	Model    domain.Model
	Instance domain.Instance
	Message  *domain.Message
}

func NewState(model domain.Model, instance domain.Instance, message *domain.Message) State {
	return State{
		Model:    model,
		Instance: instance,
		Message:  message,
	}
}

type stateJson struct {
	Model    model.ModelJson       `json:"model"`
	Instance instance.InstanceJson `json:"instance"`
	Message  *message.MessageJson  `json:"message"`
}

func (state *State) Encrypt(sender *eddsa.PrivateKey, receiver domain.PublicKey) domain.EncryptedState {
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
	return serializedState.Encrypt(sender, receiver)
}

func Decrypt(encryptedState domain.EncryptedState, receiver *eddsa.PrivateKey, sender domain.PublicKey) (State, error) {
	serializedState, err := encryptedState.Decrypt(receiver, sender)
	if err != nil {
		return State{}, err
	}
	var stateJson stateJson
	err = json.Unmarshal(serializedState.Value, &stateJson)
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
