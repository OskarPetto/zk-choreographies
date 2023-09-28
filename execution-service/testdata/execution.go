package testdata

import (
	"crypto/rand"
	"execution-service/domain"
	"execution-service/signature"
	"execution-service/utils"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

func GetPublicKeys(signatureService signature.SignatureService, count int) []domain.PublicKey {
	publicKeys := make([]domain.PublicKey, count)
	publicKeys[0] = signatureService.GetPublicKey()
	for i := 1; i < count; i++ {
		privateKey, err := eddsa.GenerateKey(rand.Reader)
		utils.PanicOnError(err)
		publicKeys[i] = domain.PublicKey{
			Value: privateKey.PublicKey.Bytes(),
		}
	}
	return publicKeys
}

func GetModel2Instance1(publicKeys []domain.PublicKey) domain.Instance {
	return getModel2Instance(
		[]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		publicKeys,
		[]domain.Hash{
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
		},
	)
}

func GetModel2Instance2(publicKeys []domain.PublicKey) domain.Instance {
	return getModel2Instance(
		[]int8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		publicKeys,
		[]domain.Hash{
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
		},
	)
}

func GetModel2Instance3(publicKeys []domain.PublicKey) domain.Instance {
	return getModel2Instance(
		[]int8{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		publicKeys,
		[]domain.Hash{
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.HashMessage([]byte("hello")),
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
		},
	)
}

func GetModel2Instance4(publicKeys []domain.PublicKey) domain.Instance {
	return getModel2Instance(
		[]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		publicKeys,
		[]domain.Hash{
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
			domain.HashMessage([]byte("hello2")),
			domain.Hash{},
			domain.HashMessage([]byte("hello3")),
			domain.Hash{},
			domain.Hash{},
			domain.HashMessage([]byte("hello")),
			domain.Hash{},
			domain.Hash{},
			domain.Hash{},
		},
	)
}

func getModel2Instance(tokenCounts []int8, publicKeys []domain.PublicKey, messageHashes []domain.Hash) domain.Instance {
	var tokenCountsFixedSize [domain.MaxPlaceCount]int8
	copy(tokenCountsFixedSize[:], tokenCounts)
	for i := len(tokenCounts); i < domain.MaxPlaceCount; i++ {
		tokenCountsFixedSize[i] = domain.InvalidTokenCount
	}
	var publicKeysFixedSize [domain.MaxParticipantCount]domain.PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := len(publicKeys); i < domain.MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = domain.InvalidPublicKey()
	}
	var messageHashesFixedSize [domain.MaxMessageCount]domain.Hash
	copy(messageHashesFixedSize[:], messageHashes)
	for i := len(messageHashes); i < domain.MaxMessageCount; i++ {
		messageHashesFixedSize[i] = domain.InvalidHash()
	}
	instance := domain.Instance{
		Model:         "example_choreography",
		TokenCounts:   tokenCountsFixedSize,
		PublicKeys:    publicKeysFixedSize,
		MessageHashes: messageHashesFixedSize,
		UpdatedAt:     time.Now().Unix(),
	}
	instance.ComputeHash()
	return instance
}
