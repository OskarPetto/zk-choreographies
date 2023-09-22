package instance

import (
	"encoding/hex"
	"fmt"
)

type InstanceService struct {
	isLoaded  bool
	instances map[string]Instance
}

var instanceService InstanceService

func NewInstanceService() InstanceService {
	if !instanceService.isLoaded {
		instanceService = InstanceService{
			isLoaded:  true,
			instances: make(map[string]Instance),
		}
	}
	return instanceService
}

func (service *InstanceService) SaveInstance(instance Instance) {
	id := hex.EncodeToString(instance.Hash)
	service.instances[id] = instance
}

func (service *InstanceService) FindInstanceByHash(hash []byte) (Instance, error) {
	id := hex.EncodeToString(hash)
	instance, exists := service.instances[id]
	if !exists {
		return Instance{}, fmt.Errorf("instance %s not found", id)
	}
	return instance, nil
}
