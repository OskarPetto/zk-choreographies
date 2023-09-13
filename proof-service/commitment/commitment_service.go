package commitment

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
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

func (service *CommitmentService) CreateCommitment(commitmentId CommitmentId, data []byte) (Commitment, error) {
	randomness := make([]byte, RandomnessSize)
	rand.Read(randomness)
	input := append(data, randomness...)
	hash := sha256.Sum256(input)
	return Commitment{
		Id:         commitmentId,
		Value:      hash,
		Randomness: randomness,
	}, nil
}
