package instance

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/message"
	"execution-service/utils"
	"fmt"
	"time"
)

type InstanceJson struct {
	Id            string        `json:"id"`
	Hash          hash.HashJson `json:"hash"`
	Model         string        `json:"model"`
	TokenCounts   []int         `json:"tokenCounts"`
	PublicKeys    []string      `json:"publicKeys"`
	MessageHashes []string      `json:"messageHashes"`
	CreatedAt     time.Time     `json:"updatedAt"`
}

type InstantiateModelCommandJson struct {
	Model      string   `json:"model"`
	PublicKeys []string `json:"publicKeys"`
}

type ExecuteTransitionCommandJson struct {
	Model                string                            `json:"model"`
	Instance             string                            `json:"instance"`
	Transition           string                            `json:"transition"`
	CreateMessageCommand *message.CreateMessageCommandJson `json:"createMessageCommand,omitempty"`
}

func ToJson(instance domain.Instance) InstanceJson {
	tokenCounts := make([]int, len(instance.TokenCounts))
	for i, tokenCount := range instance.TokenCounts {
		tokenCounts[i] = int(tokenCount)
	}
	publicKeys := make([]string, len(instance.PublicKeys))
	for i, publicKey := range instance.PublicKeys {
		publicKeys[i] = utils.BytesToString(publicKey.Value)
	}
	messageHashes := make([]string, len(instance.MessageHashes))
	for i, messageHash := range instance.MessageHashes {
		messageHashes[i] = utils.BytesToString(messageHash[:])
	}
	return InstanceJson{
		Id:            instance.Id(),
		Hash:          hash.HashToJson(instance.Hash),
		Model:         instance.Model,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
		CreatedAt:     time.Unix(instance.CreatedAt, 0),
	}
}

func (json *InstanceJson) ToInstance() (domain.Instance, error) {
	if len(json.TokenCounts) > domain.MaxPlaceCount {
		return domain.Instance{}, fmt.Errorf("instance %s has too many places", json.Id)
	}
	if len(json.PublicKeys) > domain.MaxParticipantCount {
		return domain.Instance{}, fmt.Errorf("instance %s has too many participants", json.Id)
	}
	if len(json.MessageHashes) > domain.MaxMessageCount {
		return domain.Instance{}, fmt.Errorf("instance %s has too many messages", json.Id)
	}
	tokenCounts := make([]int8, len(json.TokenCounts))
	for i, tokenCount := range json.TokenCounts {
		if tokenCount != 0 && tokenCount != 1 {
			return domain.Instance{}, fmt.Errorf("instance %s has invalid tokenCount", json.Id)
		}
		tokenCounts[i] = int8(tokenCount)
	}
	publicKeys := make([]domain.PublicKey, len(json.PublicKeys))
	for i, publicKey := range json.PublicKeys {
		publicKeyBytes, err := utils.StringToBytes(publicKey)
		if err != nil {
			return domain.Instance{}, fmt.Errorf("instance %s has invalid publicKey", json.Id)
		}
		publicKeys[i] = domain.PublicKey{
			Value: publicKeyBytes,
		}
	}
	messageHashes := make([][domain.HashSize]byte, len(json.MessageHashes))
	for i, messageHash := range json.MessageHashes {
		hash, err := utils.StringToBytes(messageHash)
		if err != nil {
			return domain.Instance{}, fmt.Errorf("instance %s has invalid messageHash", json.Id)
		}
		messageHashes[i] = [domain.HashSize]byte(hash)
	}
	hash, err := json.Hash.ToHash()
	if err != nil {
		return domain.Instance{}, fmt.Errorf("instance %s has invalid hash", json.Id)
	}
	return domain.Instance{
		Hash:          hash,
		Model:         json.Model,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
		CreatedAt:     json.CreatedAt.Unix(),
	}, nil
}

func (cmd *InstantiateModelCommandJson) ToExecutionCommand() (InstantiateModelCommand, error) {
	publicKeys := make([]domain.PublicKey, len(cmd.PublicKeys))
	for i, publicKey := range cmd.PublicKeys {
		bytes, err := utils.StringToBytes(publicKey)
		if err != nil {
			return InstantiateModelCommand{}, err
		}
		publicKeys[i] = domain.PublicKey{
			Value: bytes,
		}
	}
	return InstantiateModelCommand{
		Model:      cmd.Model,
		PublicKeys: publicKeys,
	}, nil
}

func (cmd *ExecuteTransitionCommandJson) ToExecutionCommand() (ExecuteTransitionCommand, error) {
	var createMessageCommand *message.CreateMessageCommand = nil
	if cmd.CreateMessageCommand != nil {
		result := cmd.CreateMessageCommand.ToMessageCommand()
		createMessageCommand = &result
	}
	return ExecuteTransitionCommand{
		Model:                cmd.Model,
		Instance:             cmd.Instance,
		Transition:           cmd.Transition,
		CreateMessageCommand: createMessageCommand,
	}, nil
}
