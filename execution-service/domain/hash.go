package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"execution-service/utils"
	"hash"

	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	ghash "github.com/consensys/gnark-crypto/hash"
)

const HashSize = 32
const SaltSize = 32

type HashId = string

type Hash struct {
	Value [HashSize]byte
	Salt  [SaltSize]byte
}

func EmptyHash() Hash {
	return Hash{}
}

func OutOfBoundsHash() Hash {
	var hash = Hash{}
	hash.Value[HashSize-1] = 1
	return hash
}

func (hash *Hash) String() HashId {
	return utils.BytesToString(hash.Value[:])
}

func (publicKey *PublicKey) ComputeHash() [HashSize]byte {
	var eddsaPublicKey eddsa.PublicKey
	eddsaPublicKey.A.SetBytes(publicKey.Value)
	xBytes := eddsaPublicKey.A.X.Bytes()
	yBytes := eddsaPublicKey.A.Y.Bytes()
	mimc := mimc.NewMiMC()
	mimc.Write(xBytes[:])
	mimc.Write(yBytes[:])
	return computeHash(mimc)
}

func (transition *Transition) ComputeHash() [HashSize]byte {
	mimc := mimc.NewMiMC()
	for _, incomingPlace := range transition.IncomingPlaces {
		hashUint8(mimc, incomingPlace)
	}
	for i := len(transition.IncomingPlaces); i < MaxBranchingFactor; i++ {
		hashUint8(mimc, OutOfBoundsPlaceId)
	}
	for _, outgoingPlace := range transition.OutgoingPlaces {
		hashUint8(mimc, outgoingPlace)
	}
	for i := len(transition.OutgoingPlaces); i < MaxBranchingFactor; i++ {
		hashUint8(mimc, OutOfBoundsPlaceId)
	}
	hashUint8(mimc, transition.Participant)
	hashUint8(mimc, transition.Message)
	for _, coefficient := range transition.Constraint.Coefficients {
		hashInt64(mimc, int64(coefficient))
	}
	for i := len(transition.Constraint.Coefficients); i < MaxConstraintMessageCount; i++ {
		hashInt64(mimc, 0)
	}
	for _, messageId := range transition.Constraint.MessageIds {
		hashUint8(mimc, messageId)
	}
	for i := len(transition.Constraint.MessageIds); i < MaxConstraintMessageCount; i++ {
		hashInt64(mimc, int64(EmptyMessageId))
	}
	hashInt64(mimc, int64(transition.Constraint.Offset))
	hashUint8(mimc, transition.Constraint.ComparisonOperator)
	return computeHash(mimc)
}

func (model *Model) ComputeHash() {
	mimc := mimc.NewMiMC()
	hashUint8(mimc, model.PlaceCount)
	hashUint8(mimc, model.ParticipantCount)
	hashUint8(mimc, model.MessageCount)
	for _, startPlace := range model.StartPlaces {
		hashUint8(mimc, startPlace)
	}
	for i := len(model.StartPlaces); i < MaxStartPlaceCount; i++ {
		hashUint8(mimc, OutOfBoundsPlaceId)
	}
	endPlaceTree := merkletree.New(ghash.MIMC_BN254.New())
	for _, endPlace := range model.EndPlaces {
		bytes := Uint8ToBytes(endPlace)
		endPlaceTree.Push(bytes[:])
	}
	for i := len(model.EndPlaces); i < MaxEndPlaceCount; i++ {
		bytes := Uint8ToBytes(OutOfBoundsPlaceId)
		endPlaceTree.Push(bytes[:])
	}
	mimc.Write(endPlaceTree.Root())
	transitionTree := merkletree.New(ghash.MIMC_BN254.New())
	for _, transition := range model.Transitions {
		hash := transition.ComputeHash()
		transitionTree.Push(hash[:])
	}
	for i := len(model.Transitions); i < MaxTransitionCount; i++ {
		transition := OutOfBoundsTransition()
		hash := transition.ComputeHash()
		transitionTree.Push(hash[:])
	}
	mimc.Write(transitionTree.Root())
	salt := randomFieldElement()
	mimc.Write(salt[:])
	model.Hash = Hash{
		Value: computeHash(mimc),
		Salt:  salt,
	}
}

func (instance *Instance) ComputeHash() {
	mimc := mimc.NewMiMC()
	for _, tokenCount := range instance.TokenCounts {
		hashInt64(mimc, int64(tokenCount))
	}
	for i := len(instance.TokenCounts); i < MaxPlaceCount; i++ {
		hashInt64(mimc, int64(OutOfBoundsTokenCount))
	}
	tree := merkletree.New(ghash.MIMC_BN254.New())
	for _, publicKey := range instance.PublicKeys {
		hash := publicKey.ComputeHash()
		tree.Push(hash[:])
	}
	for i := len(instance.PublicKeys); i < MaxParticipantCount; i++ {
		publicKey := OutOfBoundsPublicKey()
		hash := publicKey.ComputeHash()
		tree.Push(hash[:])
	}
	mimc.Write(tree.Root())
	for _, messageHash := range instance.MessageHashes {
		mimc.Write(messageHash[:])
	}
	for i := len(instance.MessageHashes); i < MaxMessageCount; i++ {
		hash := OutOfBoundsHash().Value
		mimc.Write(hash[:])
	}
	salt := randomFieldElement()
	mimc.Write(salt[:])
	instance.Hash = Hash{
		Value: computeHash(mimc),
		Salt:  salt,
	}
}

func (message *Message) ComputeHash() {
	var hash Hash
	if message.IsBytesMessage() {
		hash = hashBytesMessage(message.BytesMessage)
	} else {
		hash = hashIntegerMessage(message.IntegerMessage)
	}
	message.Hash = hash
}

func Uint8ToBytes(value uint8) [fr.Bytes]byte {
	var bytes [fr.Bytes]byte
	bytes[fr.Bytes-1] = value // big endian
	return bytes
}

func hashBytesMessage(message []byte) Hash {
	salt := randomFrSizedBytes()
	input := append(message, salt[:]...)
	bytesHash := sha256.Sum256(input)
	return Hash{
		Value: hashToField(bytesHash),
		Salt:  salt,
	}
}

func hashIntegerMessage(message IntegerType) Hash {
	mimc := mimc.NewMiMC()
	hashInt64(mimc, int64(message))
	salt := randomFieldElement()
	mimc.Write(salt[:])
	return Hash{
		Value: computeHash(mimc),
		Salt:  salt,
	}
}

func computeHash(hasher hash.Hash) [HashSize]byte {
	return [HashSize]byte(hasher.Sum([]byte{}))
}

func hashInt64(hasher hash.Hash, value int64) {
	var fieldElement fr.Element
	fieldElement.SetInt64(value)
	bytes := fieldElement.Bytes()
	hasher.Write(bytes[:])
}

func hashUint8(hasher hash.Hash, value uint8) {
	bytes := Uint8ToBytes(value)
	hasher.Write(bytes[:])
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
