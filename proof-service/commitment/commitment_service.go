package commitment

import (
	"crypto/rand"
	"fmt"

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

func (service *CommitmentService) CreateCommitment(commitmentId CommitmentId, data []byte) (Commitment, error) {

	hasher := mimc.NewMiMC()
	for i := range data {
		bytes := make([]byte, hasher.BlockSize())
		bytes[hasher.BlockSize()-1] = data[i] // big endian
		hasher.Write(bytes)
	}

	randomBytes := make([]byte, RandomnessSize)
	rand.Read(randomBytes)
	fieldElements, err := fr.Hash(randomBytes, []byte("CreateCommitment"), 1)
	if err != nil {
		panic(err)
	}
	randomness := fieldElements[0].Bytes()
	hasher.Write(randomness[:])
	hash := hasher.Sum([]byte{})
	return Commitment{
		Id:         commitmentId,
		Value:      hash,
		Randomness: randomness,
	}, nil
}
