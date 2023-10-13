package state

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/utils"
)

func SerializeModel(domainModel domain.Model) domain.Plaintext {
	modelJson := model.ToJson(domainModel)
	modelJsonBytes, err := json.Marshal(modelJson)
	utils.PanicOnError(err)
	plainText := domain.Plaintext{
		Value: modelJsonBytes,
	}
	return plainText
}

func SerializeInstance(domainInstance domain.Instance) domain.Plaintext {
	instanceJson := instance.ToJson(domainInstance)
	instanceJsonBytes, err := json.Marshal(instanceJson)
	utils.PanicOnError(err)
	plainText := domain.Plaintext{
		Value: instanceJsonBytes,
	}
	return plainText
}

func SerializeMessage(domainMessage domain.Message) domain.Plaintext {
	messageJson := message.ToJson(domainMessage)
	messageJsonBytes, err := json.Marshal(messageJson)
	utils.PanicOnError(err)
	plainText := domain.Plaintext{
		Value: messageJsonBytes,
	}
	return plainText
}

func DeserializeModel(plaintext domain.Plaintext) (domain.Model, error) {
	var modelJson model.ModelJson
	err := json.Unmarshal(plaintext.Value, &modelJson)
	if err != nil {
		return domain.Model{}, err
	}
	model, err := modelJson.ToModel()
	if err != nil {
		return domain.Model{}, err
	}
	return model, nil
}

func DeserializeInstance(plaintext domain.Plaintext) (domain.Instance, error) {
	var instanceJson instance.InstanceJson
	err := json.Unmarshal(plaintext.Value, &instanceJson)
	if err != nil {
		return domain.Instance{}, err
	}
	instance, err := instanceJson.ToInstance()
	if err != nil {
		return domain.Instance{}, err
	}
	return instance, nil
}

func DeserializeMessage(plaintext domain.Plaintext) (domain.Message, error) {
	var messageJson message.MessageJson
	err := json.Unmarshal(plaintext.Value, &messageJson)
	if err != nil {
		return domain.Message{}, err
	}
	message, err := messageJson.ToMessage()
	if err != nil {
		return domain.Message{}, err
	}
	return message, nil
}
