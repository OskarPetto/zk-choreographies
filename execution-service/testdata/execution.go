package testdata

import (
	"crypto/rand"
	"proof-service/crypto"
	"proof-service/domain"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

var signatureService crypto.SignatureService = crypto.NewSignatureService()

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

func GetPetriNet1Instance1(publicKey []byte) domain.Instance {
	return domain.Instance{
		Id:          "conformance_example1",
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetPetriNet1Instance2(publicKey []byte) domain.Instance {
	return domain.Instance{
		Id:          "conformance_example2",
		TokenCounts: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetPetriNet1Instance3(publicKey []byte) domain.Instance {
	return domain.Instance{
		Id:          "conformance_example4",
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetCommitment1() crypto.SaltedHash {
	return crypto.SaltedHash{
		Value: []byte{15, 119, 210, 82, 4, 149, 235, 173, 255, 201, 90, 205, 146, 233, 251, 58, 54, 88, 10, 179, 75, 101, 147, 46, 127, 239, 221, 252, 28, 71, 138, 66},
		Salt:  []byte{85, 39, 212, 198, 200, 84, 236, 218, 89, 123, 119, 127, 251, 16, 159, 125, 24, 72, 146, 14, 13, 242, 101, 182, 18, 14, 139, 149, 217, 116, 255, 43},
	}
}
