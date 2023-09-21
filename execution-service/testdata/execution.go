package testdata

import (
	"proof-service/crypto"
	"proof-service/execution"
)

func GetPetriNet1Instance1(publicKey []byte) execution.Instance {
	return execution.Instance{
		Id:          "conformance_example1",
		TokenCounts: []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetPetriNet1Instance2(publicKey []byte) execution.Instance {
	return execution.Instance{
		Id:          "conformance_example2",
		TokenCounts: []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		PublicKeys: [][]byte{
			publicKey,
		},
	}
}

func GetPetriNet1Instance3(publicKey []byte) execution.Instance {
	return execution.Instance{
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

func GetPublicKey1() []byte {
	return []byte{70, 200, 160, 220, 129, 215, 38, 174, 106, 10, 190, 160, 109, 87, 219, 147, 161, 184, 34, 209, 190, 54, 152, 202, 123, 230, 254, 52, 193, 43, 56, 147}
}
