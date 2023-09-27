package json

import (
	"bytes"
	"encoding/base64"
	"execution-service/domain"
	"time"
)

type Hash struct {
	Value string `json:"value"`
	Salt  string `json:"salt"`
}

type Instance struct {
	Id            string    `json:"id"`
	Hash          Hash      `json:"hash"`
	Model         string    `json:"model"`
	TokenCounts   []int     `json:"tokenCounts"`
	PublicKeys    []string  `json:"publicKeys"`
	MessageHashes []Hash    `json:"messageHashes"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func FromDomainInstance(instance domain.Instance) Instance {
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
	messageHashes := make([]Hash, 0)
	for _, messageHash := range instance.MessageHashes {
		invalidHash := domain.InvalidHash()
		if bytes.Equal(invalidHash.Value[:], messageHash.Value[:]) && bytes.Equal(invalidHash.Salt[:], messageHash.Salt[:]) {
			break
		}
		messageHashes = append(messageHashes, Hash{
			Value: bytesToString(messageHash.Value[:]),
			Salt:  bytesToString(messageHash.Salt[:]),
		})
	}
	return Instance{
		Id: instance.Id(),
		Hash: Hash{
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
