package commitment

import (
	"crypto/rand"
	"fmt"
	"proof-service/workflow"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

type CommitmentService struct {
	commitments map[CommitmentId]Commitment
}

func NewCommitmentService() CommitmentService {
	return CommitmentService{
		commitments: make(map[CommitmentId]Commitment),
	}
}

func (service *CommitmentService) FindCommitment(commitmentId CommitmentId) (Commitment, error) {
	commitment, exists := service.commitments[commitmentId]
	if !exists {
		return Commitment{}, fmt.Errorf("randomness for %s not found", commitmentId)
	}
	return commitment, nil
}

func (service *CommitmentService) CreateCommitment(instance workflow.Instance) Commitment {

	serializedInstance := serializeInstance(instance)
	randomness := randomFieldElement()
	value := commitBytes(serializedInstance, randomness)

	commitment := Commitment{
		Id:         instance.Id,
		Value:      value,
		Randomness: randomness,
	}
	return commitment
}

func (service *CommitmentService) SaveCommitment(commitment Commitment) {
	service.commitments[commitment.Id] = commitment
}

func commitBytes(input []byte, randomness [mimc.BlockSize]byte) []byte {
	hasher := mimc.NewMiMC()
	for i := range input {
		bytes := make([]byte, hasher.BlockSize())
		bytes[hasher.BlockSize()-1] = input[i] // big endian
		hasher.Write(bytes)
	}
	hasher.Write(randomness[:])
	return hasher.Sum([]byte{})
}

func serializeInstance(instance workflow.Instance) []byte {
	placeCount := len(instance.TokenCounts)
	var bytes = make([]byte, workflow.MaxPlaceCount+1)
	bytes[0] = byte(placeCount)
	for i := 0; i < placeCount; i++ {
		bytes[i+1] = byte(instance.TokenCounts[i])
	}
	return bytes
}

func randomFieldElement() [mimc.BlockSize]byte {
	randomBytes := make([]byte, mimc.BlockSize)
	rand.Read(randomBytes)
	fieldElements, err := fr.Hash(randomBytes, []byte("randomFieldElement"), 1)
	if err != nil {
		panic(err)
	}
	fieldElementBytes := fieldElements[0].Bytes()
	return fieldElementBytes
}
