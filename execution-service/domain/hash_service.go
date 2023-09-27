package domain

import "fmt"

type HashService struct {
	isLoaded bool
	hashes   map[string]Hash
}

var hashService HashService

func NewHashService() HashService {
	if !hashService.isLoaded {
		hashService = HashService{
			isLoaded: true,
			hashes:   make(map[string]Hash),
		}
	}
	return hashService
}

func (service *HashService) SaveModelHash(modelId ModelId, hash Hash) {
	service.hashes[modelId] = hash
}

func (service *HashService) FindHashByModelId(id ModelId) (Hash, error) {
	hash, exists := service.hashes[id]
	if !exists {
		return Hash{}, fmt.Errorf("hash for model %s not found", id)
	}
	return hash, nil
}
