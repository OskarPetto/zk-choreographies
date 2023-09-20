package crypto

import (
	"proof-service/utils"
	"proof-service/workflow"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

type CommitmentId = string

type Commitment struct {
	Value      []byte
	Randomness []byte
}

func Commit(instance workflow.Instance) Commitment {
	serializedInstance := serializeInstance(instance)
	fieldElements := bytesToFieldElements(serializedInstance)
	randomness := randomFieldElement()
	input := append(fieldElements, randomness...)
	value := hashFieldElements(input)

	commitment := Commitment{
		Value:      value,
		Randomness: randomness,
	}
	return commitment
}

func hashFieldElements(input []byte) []byte {
	res, err := mimc.Sum(input)
	utils.PanicOnError(err)
	return res
}

func serializeInstance(instance workflow.Instance) []byte {
	var bytes = make([]byte, workflow.MaxPlaceCount+1)
	placeCount := len(instance.TokenCounts)
	bytes[0] = byte(placeCount)
	for i := 0; i < placeCount; i++ {
		bytes[i+1] = byte(instance.TokenCounts[i])
	}
	return bytes
}
