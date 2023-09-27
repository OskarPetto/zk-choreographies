package testdata

import (
	"crypto/rand"
	"crypto/sha256"
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
		MessageHashes: [domain.MaxMessageCount]domain.MessageHash{
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
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
		MessageHashes: [domain.MaxMessageCount]domain.MessageHash{
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
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
		MessageHashes: [domain.MaxMessageCount]domain.MessageHash{
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.MessageHash{
				Value: sha256.Sum256([]byte("hello")),
			},
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
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
		MessageHashes: [domain.MaxMessageCount]domain.MessageHash{
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
			domain.DefaultMessageHash,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
	}
}
