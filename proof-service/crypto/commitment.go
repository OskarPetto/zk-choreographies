package crypto

import (
	"crypto/rand"
	"proof-service/utils"
	"proof-service/workflow"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

type CommitmentId = string

type Commitment struct {
	Value      []byte
	Randomness []byte
}

func Commit(instance workflow.Instance) Commitment {
	fieldElements := serializeInstance(instance)
	randomness := randomFieldElement()
	input := append(fieldElements, randomness...)

	value, err := mimc.Sum(input)
	utils.PanicOnError(err)

	commitment := Commitment{
		Value:      value,
		Randomness: randomness,
	}
	return commitment
}

func serializeInstance(instance workflow.Instance) []byte {
	var bytes = make([]byte, (workflow.MaxPlaceCount+1)*fr.Bytes)
	placeCount := len(instance.TokenCounts)
	bytes[fr.Bytes-1] = byte(placeCount) // big endian
	for i := 0; i < placeCount; i++ {
		bytes[(i+1)*fr.Bytes+fr.Bytes-1] = byte(instance.TokenCounts[i]) // big endian
	}
	return bytes
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
