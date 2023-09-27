package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"math/big"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

func hashMessage(message []byte) MessageHash {
	return MessageHash{
		Value: sha256.Sum256(message),
	}
}

func (model *Model) ComputeHash() {
	mimc := mimc.NewMiMC()
	hashUint8(mimc, model.PlaceCount)
	hashUint8(mimc, model.ParticipantCount)
	hashUint8(mimc, model.MessageCount)
	for _, startPlace := range model.StartPlaces {
		hashUint8(mimc, startPlace)
	}
	for _, endPlace := range model.EndPlaces {
		hashUint8(mimc, endPlace)
	}
	for _, transition := range model.Transitions {
		hashUint8(mimc, boolToUint8(transition.IsInitialized))
		for _, incomingPlace := range transition.IncomingPlaces {
			hashUint8(mimc, incomingPlace)
		}
		for _, outgoingPlace := range transition.OutgoingPlaces {
			hashUint8(mimc, outgoingPlace)
		}
		hashUint8(mimc, transition.Participant)
		hashUint8(mimc, transition.Message)
	}
	model.Salt = randomFieldElement()
	mimc.Write(model.Salt)
	model.Hash = mimc.Sum([]byte{})
}

func (instance *Instance) ComputeHash() {
	mimc := mimc.NewMiMC()
	for _, tokenCount := range instance.TokenCounts {
		hashInt8(mimc, tokenCount)
	}
	for _, publicKey := range instance.PublicKeys {
		var eddsaPublicKey eddsa.PublicKey
		eddsaPublicKey.A.SetBytes(publicKey.Value)
		xBytes := eddsaPublicKey.A.X.Bytes()
		yBytes := eddsaPublicKey.A.Y.Bytes()
		mimc.Write(xBytes[:])
		mimc.Write(yBytes[:])
	}
	for _, messageHash := range instance.MessageHashes {
		fieldElement := utils.HashToField(messageHash.Value)
		mimc.Write(fieldElement)
	}
	instance.Salt = randomFieldElement()
	mimc.Write(instance.Salt)
	instance.Hash = mimc.Sum([]byte{})
}

func hashInt8(hasher hash.Hash, value int8) {
	var bytes [fr.Bytes]byte
	number := big.NewInt(int64(value))
	number.FillBytes(bytes[:])
	hasher.Write(bytes[:])
}

func hashUint8(hasher hash.Hash, value uint8) {
	var bytes [fr.Bytes]byte
	bytes[fr.Bytes-1] = value // big endian
	hasher.Write(bytes[:])
}

func boolToUint8(value bool) uint8 {
	var result uint8 = 0
	if value {
		result = 1
	}
	return result
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
