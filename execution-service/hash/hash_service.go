package hash

import (
	"execution-service/domain"
	"fmt"
)

type HashService struct {
	hashes map[string]domain.Hash
}

func NewHashService() HashService {
	return HashService{
		hashes: make(map[string]domain.Hash),
	}
}

func (service *HashService) SaveModelHash(modelId domain.ModelId, hash domain.Hash) {
	service.hashes[modelId] = hash
}

func (service *HashService) FindHashByModelId(id domain.ModelId) (domain.Hash, error) {
	hash, exists := service.hashes[id]
	if !exists {
		return domain.Hash{}, fmt.Errorf("hash for model %s not found", id)
	}
	return hash, nil
}
