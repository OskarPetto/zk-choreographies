package crypto

import (
	"crypto/rand"
	"proof-service/domain"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

type SaltedHashId = string

type SaltedHash struct {
	Value []byte
	Salt  []byte
}

func HashInstance(instance domain.Instance) SaltedHash {
	salt := randomFieldElement()
	return SaltedHash{
		Value: hashInstance(instance, salt),
		Salt:  salt,
	}
}

func hashInstance(instance domain.Instance, salt []byte) []byte {
	mimc := mimc.NewMiMC()
	for _, tokenCount := range instance.TokenCounts {
		var bytes [fr.Bytes]byte
		bytes[fr.Bytes-1] = byte(tokenCount) // big endian
		mimc.Write(bytes[:])
	}
	for i := len(instance.TokenCounts); i < domain.MaxPlaceCount; i++ {
		var bytes [fr.Bytes]byte
		mimc.Write(bytes[:])
	}
	for _, publicKeyBytes := range instance.PublicKeys {
		var publicKey eddsa.PublicKey
		publicKey.A.SetBytes(publicKeyBytes)
		xBytes := publicKey.A.X.Bytes()
		yBytes := publicKey.A.Y.Bytes()
		mimc.Write(xBytes[:])
		mimc.Write(yBytes[:])
	}
	for i := len(instance.PublicKeys); i < domain.MaxParticipantCount; i++ {
		var zeros [fr.Bytes]byte
		mimc.Write(zeros[:])
		mimc.Write(zeros[:])
	}
	mimc.Write(salt)
	return mimc.Sum([]byte{})
}

func randomFieldElement() []byte {
	randomBytes := randomFrSizedBytes()
	fieldElements, err := fr.Hash(randomBytes, []byte("randomFieldElement"), 1)
	utils.PanicOnError(err)
	fieldElementBytes := fieldElements[0].Bytes()
	return fieldElementBytes[:]
}

func randomFrSizedBytes() []byte {
	res := make([]byte, fr.Bytes)
	_, err := rand.Read(res)
	utils.PanicOnError(err)
	return res
}
