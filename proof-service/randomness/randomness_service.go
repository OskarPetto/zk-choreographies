package randomness

import (
	"crypto/rand"
	"fmt"
	"proof-service/workflow"
)

const byteCound = 32

type Randomness = []byte

var defaultRandomness = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type RandomnessService struct {
	Randomnesses map[workflow.InstanceId]Randomness
}

func NewRandomnessService() RandomnessService {
	return RandomnessService{
		Randomnesses: make(map[workflow.InstanceId]Randomness),
	}
}

func (service *RandomnessService) FindRandomness(instanceId workflow.InstanceId) (Randomness, error) {
	randomness, exists := service.Randomnesses[instanceId]
	if !exists {
		return defaultRandomness, fmt.Errorf("Randomness for %s not found", instanceId)
	}
	return randomness, nil
}

func (service *RandomnessService) CreateRandomness(instanceId workflow.InstanceId) (Randomness, error) {
	newRandomness := make([]byte, byteCound)
	rand.Read(newRandomness)
	service.Randomnesses[instanceId] = newRandomness
	return newRandomness, nil
}
