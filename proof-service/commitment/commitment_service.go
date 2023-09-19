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

func (service *CommitmentService) SaveCommitment(commitment Commitment) error {
	service.commitments[commitment.Id] = commitment
	return nil
}

func (service *CommitmentService) CreateCommitment(instance workflow.Instance) (Commitment, error) {

	serializedInstance := serializeInstance(instance)
	hasher := mimc.NewMiMC()
	for i := range serializedInstance {
		bytes := make([]byte, hasher.BlockSize())
		bytes[hasher.BlockSize()-1] = serializedInstance[i] // big endian
		hasher.Write(bytes)
	}

	randomness := randomFieldElement()
	hasher.Write(randomness[:])
	hash := hasher.Sum([]byte{})
	return Commitment{
		Id:         instance.Id,
		Value:      hash,
		Randomness: randomness,
	}, nil
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
