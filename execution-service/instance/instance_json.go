package instance

import (
	"bytes"
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/utils"
	"fmt"
	"time"
)

type InstanceJson struct {
	Id            string          `json:"id"`
	Hash          hash.HashJson   `json:"hash"`
	Model         string          `json:"model"`
	TokenCounts   []int           `json:"tokenCounts"`
	PublicKeys    []string        `json:"publicKeys"`
	MessageHashes []hash.HashJson `json:"messageHashes"`
	CreatedAt     time.Time       `json:"updatedAt"`
}

func ToJson(instance domain.Instance) InstanceJson {
	tokenCounts := make([]int, 0)
	for _, tokenCount := range instance.TokenCounts {
		if tokenCount == domain.InvalidTokenCount {
			break
		}
		tokenCounts = append(tokenCounts, int(tokenCount))
	}
	publicKeys := make([]string, 0)
	for _, publicKey := range instance.PublicKeys {
		if bytes.Equal(domain.InvalidPublicKey().Value, publicKey.Value) {
			break
		}
		publicKeys = append(publicKeys, utils.BytesToString(publicKey.Value))
	}
	messageHashes := make([]hash.HashJson, 0)
	for _, messageHash := range instance.MessageHashes {
		invalidHash := domain.InvalidHash()
		if bytes.Equal(invalidHash.Value[:], messageHash.Value[:]) && bytes.Equal(invalidHash.Salt[:], messageHash.Salt[:]) {
			break
		}
		messageHashes = append(messageHashes, hash.HashToJson(messageHash))
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
	var tokenCounts [domain.MaxPlaceCount]int8
	for i, tokenCount := range json.TokenCounts {
		if tokenCount != 0 && tokenCount != 1 {
			return domain.Instance{}, fmt.Errorf("instance %s has invalid tokenCount", json.Id)
		}
		tokenCounts[i] = int8(tokenCount)
	}
	for i := len(json.TokenCounts); i < domain.MaxPlaceCount; i++ {
		tokenCounts[i] = domain.InvalidTokenCount
	}

	var publicKeys [domain.MaxParticipantCount]domain.PublicKey
	for i, publicKey := range json.PublicKeys {
		publicKeyBytes, err := utils.StringToBytes(publicKey)
		if err != nil {
			return domain.Instance{}, fmt.Errorf("instance %s has invalid publicKey", json.Id)
		}
		publicKeys[i] = domain.PublicKey{
			Value: publicKeyBytes,
		}
	}
	for i := len(json.PublicKeys); i < domain.MaxParticipantCount; i++ {
		publicKeys[i] = domain.InvalidPublicKey()
	}

	var messageHashes [domain.MaxMessageCount]domain.Hash
	for i, messageHash := range json.MessageHashes {
		hash, err := messageHash.ToHash()
		if err != nil {
			return domain.Instance{}, fmt.Errorf("instance %s has invalid messageHash", json.Id)
		}
		messageHashes[i] = hash
	}
	for i := len(json.MessageHashes); i < domain.MaxMessageCount; i++ {
		messageHashes[i] = domain.InvalidHash()
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
