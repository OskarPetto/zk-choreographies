package crypto

import (
	"crypto/rand"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func bytesToFieldElements(bytes []byte) []byte {
	var fieldElements []byte
	for i := range bytes {
		fieldElement := make([]byte, fr.Bytes)
		fieldElement[fr.Bytes-1] = bytes[i] // big endian
		fieldElements = append(fieldElements, fieldElement...)
	}
	return fieldElements
}

func randomFieldElement() []byte {
	randomBytes := randomFrSizedBytes()
	return hashToField(randomBytes)
}

func randomFrSizedBytes() []byte {
	res := make([]byte, fr.Bytes)
	_, err := rand.Read(res)
	utils.PanicOnError(err)
	return res
}

func hashToField(data []byte) []byte {
	fieldElements, err := fr.Hash(data, []byte("randomFieldElement"), 1)
	utils.PanicOnError(err)
	fieldElementBytes := fieldElements[0].Bytes()
	return fieldElementBytes[:]
}
