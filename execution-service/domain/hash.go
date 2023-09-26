package domain

import (
	"crypto/rand"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

func (instance *Instance) ComputeHash() {
	instance.Salt = randomFieldElement()
	mimc := mimc.NewMiMC()
	for _, tokenCount := range instance.TokenCounts {
		var bytes [fr.Bytes]byte
		bytes[fr.Bytes-1] = byte(tokenCount) // big endian
		mimc.Write(bytes[:])
	}
	for i := len(instance.TokenCounts); i < MaxPlaceCount; i++ {
		var zeros [fr.Bytes]byte
		mimc.Write(zeros[:])
	}
	for _, publicKey := range instance.PublicKeys {
		var eddsaPublicKey eddsa.PublicKey
		eddsaPublicKey.A.SetBytes(publicKey.Value)
		xBytes := eddsaPublicKey.A.X.Bytes()
		yBytes := eddsaPublicKey.A.Y.Bytes()
		mimc.Write(xBytes[:])
		mimc.Write(yBytes[:])
	}
	for i := len(instance.PublicKeys); i < MaxParticipantCount; i++ {
		var zeros [fr.Bytes]byte
		mimc.Write(zeros[:])
		mimc.Write(zeros[:])
	}
	for _, messageHash := range instance.MessageHashes {
		fieldElement := utils.HashToField(messageHash.Value)
		mimc.Write(fieldElement[:])
	}
	for i := len(instance.MessageHashes); i < MaxMessageCount; i++ {
		var zeros [fr.Bytes]byte
		fieldElement := utils.HashToField(zeros)
		mimc.Write(fieldElement[:])
	}
	mimc.Write(instance.Salt)
	instance.Hash = mimc.Sum([]byte{})
}

func randomFieldElement() []byte {
	randomBytes := randomFrSizedBytes()
	fieldElement := utils.HashToField(randomBytes)
	return fieldElement[:]
}

func randomFrSizedBytes() [fr.Bytes]byte {
	res := make([]byte, fr.Bytes)
	_, err := rand.Read(res)
	utils.PanicOnError(err)
	return [fr.Bytes]byte(res)
}
