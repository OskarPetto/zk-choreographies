package testdata

import (
	"crypto/rand"
	"proof-service/authentication"
	"proof-service/instance"
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

func GetPetriNet1Instance1(publicKey []byte) instance.Instance {
	return instance.Instance{
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
	}
}

func GetPetriNet1Instance2(publicKey []byte) instance.Instance {
	return instance.Instance{
		TokenCounts: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
	}
}

func GetPetriNet1Instance3(publicKey []byte) instance.Instance {
	return instance.Instance{
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		PublicKeys: [][]byte{
			publicKey,
		},
		Hash: make([]byte, 32),
		Salt: make([]byte, 32),
	}
}
