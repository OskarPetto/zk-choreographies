package json

import (
	"bytes"
	"encoding/hex"
	"execution-service/domain"
)

type Hash struct {
	Value string `json:"value"`
	Salt  string `json:"salt"`
}

type Instance struct {
	Id            string   `json:"id"`
	Hash          Hash     `json:"hash"`
	Model         string   `json:"model"`
	TokenCounts   []int    `json:"tokenCounts"`
	PublicKeys    []string `json:"publicKeys"`
	MessageHashes []Hash   `json:"messageHashes"`
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
		publicKeys = append(publicKeys, hex.EncodeToString(publicKey.Value))
	}
	messageHashes := make([]Hash, 0)
	for _, messageHash := range instance.MessageHashes {
		invalidHash := domain.InvalidHash()
		if bytes.Equal(invalidHash.Value[:], messageHash.Value[:]) && bytes.Equal(invalidHash.Salt[:], messageHash.Salt[:]) {
			break
		}
		messageHashes = append(messageHashes, Hash{
			Value: hex.EncodeToString(messageHash.Value[:]),
			Salt:  hex.EncodeToString(messageHash.Salt[:]),
		})
	}
	return Instance{
		Id: instance.Id,
		Hash: Hash{
			Value: hex.EncodeToString(instance.Hash.Value[:]),
			Salt:  hex.EncodeToString(instance.Hash.Salt[:]),
		},
		Model:         instance.Model,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
	}
}
