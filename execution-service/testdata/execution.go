package testdata

import (
	"crypto/rand"
	"proof-service/authentication"
	"proof-service/domain"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

var signatureService authentication.SignatureService = authentication.NewSignatureService()

func GetPublicKeys(count int) [][]byte {
	publicKeys := make([][]byte, count)
	publicKeys[0] = signatureService.GetPublicKey()
	for i := 1; i < count; i++ {
		privateKey, err := eddsa.GenerateKey(rand.Reader)
		utils.PanicOnError(err)
		publicKeys[i] = privateKey.PublicKey.Bytes()
	}
	return publicKeys
}

func GetModel1Instance1(publicKey []byte) domain.Instance {
	return domain.Instance{
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
	}
}

func GetModel1Instance2(publicKey []byte) domain.Instance {
	return domain.Instance{
		TokenCounts: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
	}
}

func GetModel1Instance3(publicKey []byte) domain.Instance {
	return domain.Instance{
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		PublicKeys: [][]byte{
			publicKey,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
	}
}
