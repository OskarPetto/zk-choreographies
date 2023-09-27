package testdata

import (
	"crypto/rand"
	"proof-service/authentication"
	"proof-service/domain"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

var signatureService authentication.SignatureService = authentication.NewSignatureService()

func GetPublicKeys(count int) []domain.PublicKey {
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
	var publicKeysFixedSize [domain.MaxParticipantCount]domain.PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := len(publicKeys); i < domain.MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = domain.DefaultPublicKey
	}
	return domain.Instance{
		Id:          "example_choreography1",
		Model:       "example_choreography",
		TokenCounts: [domain.MaxPlaceCount]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		PublicKeys:  publicKeysFixedSize,
		MessageHashes: [domain.MaxMessageCount]domain.Hash{
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
		},
		Hash: domain.DefaultHash,
	}
}

func GetModel2Instance2(publicKeys []domain.PublicKey) domain.Instance {
	var publicKeysFixedSize [domain.MaxParticipantCount]domain.PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := len(publicKeys); i < domain.MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = domain.DefaultPublicKey
	}
	return domain.Instance{
		Id:          "example_choreography1",
		Model:       "example_choreography",
		TokenCounts: [domain.MaxPlaceCount]int8{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		PublicKeys:  publicKeysFixedSize,
		MessageHashes: [domain.MaxMessageCount]domain.Hash{
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
		},
		Hash: domain.DefaultHash,
	}
}

func GetModel2Instance3(publicKeys []domain.PublicKey) domain.Instance {
	var publicKeysFixedSize [domain.MaxParticipantCount]domain.PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := len(publicKeys); i < domain.MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = domain.DefaultPublicKey
	}
	return domain.Instance{
		Id:          "example_choreography1",
		Model:       "example_choreography",
		TokenCounts: [domain.MaxPlaceCount]int8{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
		PublicKeys:  publicKeysFixedSize,
		MessageHashes: [domain.MaxMessageCount]domain.Hash{
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.HashMessage([]byte("hello")),
		},
		Hash: domain.DefaultHash,
	}
}

func GetModel2Instance4(publicKeys []domain.PublicKey) domain.Instance {
	var publicKeysFixedSize [domain.MaxParticipantCount]domain.PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := len(publicKeys); i < domain.MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = domain.DefaultPublicKey
	}
	return domain.Instance{
		Id:          "example_choreography1",
		Model:       "example_choreography",
		TokenCounts: [domain.MaxPlaceCount]int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		PublicKeys:  publicKeysFixedSize,
		MessageHashes: [domain.MaxMessageCount]domain.Hash{
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
			domain.DefaultHash,
		},
		Hash: domain.DefaultHash,
	}
}
