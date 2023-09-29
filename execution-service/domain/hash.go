package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"execution-service/utils"
	"hash"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

const HashSize = 32
const SaltSize = 32

type HashId = string

type Hash struct {
	Value [HashSize]byte // field element
	Salt  [SaltSize]byte
}

func EmptyHash() Hash {
	return Hash{}
}

func OutOfBoundsHash() Hash {
	var hash = Hash{}
	hash.Value[HashSize-1] = 1
	hash.Salt[SaltSize-1] = 1
	return hash
}

func HashMessage(message []byte) Hash {
	salt := randomFrSizedBytes()
	input := append(message, salt[:]...)
	bytesHash := sha256.Sum256(input)
	return Hash{
		Value: hashToField(bytesHash),
		Salt:  salt,
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
		hashUint8(mimc, boolToUint8(transition.IsValid))
		for _, incomingPlace := range transition.IncomingPlaces {
			hashUint8(mimc, incomingPlace)
		}
		for _, outgoingPlace := range transition.OutgoingPlaces {
			hashUint8(mimc, outgoingPlace)
		}
		hashUint8(mimc, transition.Participant)
		hashUint8(mimc, transition.Message)
	}
	hashInt64(mimc, model.CreatedAt)
	salt := randomFieldElement()
	mimc.Write(salt[:])
	model.Hash = Hash{
		Value: [HashSize]byte(mimc.Sum([]byte{})),
		Salt:  salt,
	}
}

func (instance *Instance) ComputeHash() {
	mimc := mimc.NewMiMC()
	for _, tokenCount := range instance.TokenCounts {
		hashInt64(mimc, int64(tokenCount))
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
		mimc.Write(messageHash.Value[:])
	}
	hashInt64(mimc, instance.CreatedAt)
	salt := randomFieldElement()
	mimc.Write(salt[:])
	instance.Hash = Hash{
		Value: [HashSize]byte(mimc.Sum([]byte{})),
		Salt:  salt,
	}
}

func hashInt64(hasher hash.Hash, value int64) {
	var fieldElement fr.Element
	fieldElement.SetInt64(value)
	bytes := fieldElement.Bytes()
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

func randomFieldElement() [fr.Bytes]byte {
	randomBytes := randomFrSizedBytes()
	fieldElement := hashToField(randomBytes)
	return fieldElement
}

func randomFrSizedBytes() [fr.Bytes]byte {
	res := make([]byte, fr.Bytes)
	_, err := rand.Read(res)
	utils.PanicOnError(err)
	return [fr.Bytes]byte(res)
}

func hashToField(data [fr.Bytes]byte) [fr.Bytes]byte {
	fieldElements, err := fr.Hash(data[:], []byte("c5f6c44a-050b-469d-8d5d-a66992a40ca7"), 1)
	utils.PanicOnError(err)
	fieldElementBytes := fieldElements[0].Bytes()
	return fieldElementBytes
}
