package message

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/utils"
)

func SerializeMessage(domainMessage domain.Message) domain.Plaintext {
	messageJson := MessageToJson(domainMessage)
	messageJsonBytes, err := json.Marshal(messageJson)
	utils.PanicOnError(err)
	plainText := domain.Plaintext{
		Value: messageJsonBytes,
	}
	return plainText
}

func DeserializeMessage(plaintext domain.Plaintext) (domain.Message, error) {
	var messageJson MessageJson
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
