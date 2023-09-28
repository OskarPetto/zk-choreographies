package instance

import (
	"bytes"
	"encoding/base64"
	"execution-service/domain"
	"time"
)

type HashJson struct {
	Value string `json:"value"`
	Salt  string `json:"salt"`
}

type InstanceJson struct {
	Id            string     `json:"id"`
	Hash          HashJson   `json:"hash"`
	Model         string     `json:"model"`
	TokenCounts   []int      `json:"tokenCounts"`
	PublicKeys    []string   `json:"publicKeys"`
	MessageHashes []HashJson `json:"messageHashes"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

func InstanceToJson(instance domain.Instance) InstanceJson {
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
		publicKeys = append(publicKeys, bytesToString(publicKey.Value))
	}
	messageHashes := make([]HashJson, 0)
	for _, messageHash := range instance.MessageHashes {
		invalidHash := domain.InvalidHash()
		if bytes.Equal(invalidHash.Value[:], messageHash.Value[:]) && bytes.Equal(invalidHash.Salt[:], messageHash.Salt[:]) {
			break
		}
		messageHashes = append(messageHashes, HashJson{
			Value: bytesToString(messageHash.Value[:]),
			Salt:  bytesToString(messageHash.Salt[:]),
		})
	}
	return InstanceJson{
		Id: instance.Id(),
		Hash: HashJson{
			Value: bytesToString(instance.Hash.Value[:]),
			Salt:  bytesToString(instance.Hash.Salt[:]),
		},
		Model:         instance.Model,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
		UpdatedAt:     time.Unix(instance.UpdatedAt, 0),
	}
}

func bytesToString(bytes []byte) string {
	return base64.StdEncoding.EncodeToString(bytes)
}
