package instance

import (
	"bytes"
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/utils"
	"time"
)

type InstanceJson struct {
	Id            string          `json:"id"`
	Hash          hash.HashJson   `json:"hash"`
	Model         string          `json:"model"`
	TokenCounts   []int           `json:"tokenCounts"`
	PublicKeys    []string        `json:"publicKeys"`
	MessageHashes []hash.HashJson `json:"messageHashes"`
	UpdatedAt     time.Time       `json:"updatedAt"`
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
		UpdatedAt:     time.Unix(instance.UpdatedAt, 0),
	}
}
