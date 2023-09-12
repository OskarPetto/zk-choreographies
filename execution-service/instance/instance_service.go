package instance

import (
	"fmt"

	"github.com/google/uuid"
)

type IInstanceService interface {
	SaveInstance(Instance) error
	FindInstance(InstanceId) (Instance, error)
}

type InstanceService struct {
	Instances map[InstanceId]Instance
}

func NewInstanceService() *InstanceService {
	return &InstanceService{
		Instances: make(map[InstanceId]Instance),
	}
}

func (service *InstanceService) SaveInstance(instance Instance) error {
	if instance.Id == "" {
		instance.Id = createInstanceId()
	}
	service.Instances[instance.Id] = instance
	return nil
}

func (service *InstanceService) FindInstance(instanceId InstanceId) (Instance, error) {
	instance, exists := service.Instances[instanceId]
	if !exists {
		return Instance{}, fmt.Errorf("Instance with %s not found", string(instanceId))
	}
	return instance, nil
}

func createInstanceId() InstanceId {
	return InstanceId(uuid.New().String())
}
