package utils

import "github.com/consensys/gnark-crypto/ecc/bn254/fr"

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func HashToField(data [fr.Bytes]byte) []byte {
	fieldElements, err := fr.Hash(data[:], []byte("HashToField"), 1)
	PanicOnError(err)
	fieldElementBytes := fieldElements[0].Bytes()
	return fieldElementBytes[:]
}
