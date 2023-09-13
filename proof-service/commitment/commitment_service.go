package commitment

import (
	"crypto/rand"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
)

type CommitmentService struct {
	Commitments map[CommitmentId]Commitment
}

func NewCommitmentService() CommitmentService {
	return CommitmentService{
		Commitments: make(map[CommitmentId]Commitment),
	}
}

func (service *CommitmentService) FindCommitment(commitmentId CommitmentId) (Commitment, error) {
	commitment, exists := service.Commitments[commitmentId]
	if !exists {
		return Commitment{}, fmt.Errorf("randomness for %s not found", commitmentId)
	}
	return commitment, nil
}

func (service *CommitmentService) SaveCommitment(commitment Commitment) error {
	service.Commitments[commitment.Id] = commitment
	return nil
}

func (service *CommitmentService) CreateCommitment(commitmentId CommitmentId, data []byte) Commitment {
	randomness := make([]byte, randomnessSize)
	rand.Read(randomness)
	input := append(data, randomness...)
	hash := mimc.NewMiMC()
	hash.Write(input)
	return Commitment{
		Id:         commitmentId,
		Value:      hash.Sum([]byte{}),
		Randomness: randomness,
	}
}
