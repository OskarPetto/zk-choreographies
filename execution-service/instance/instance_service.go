package instance

import (
	"execution-service/model"
	"fmt"

	"github.com/google/uuid"
)

type InstanceService struct {
	Instances map[InstanceId]Instance
}

func NewInstanceService() *InstanceService {
	return &InstanceService{
		Instances: make(map[InstanceId]Instance),
	}
}

func (service *InstanceService) InstantiateModel(model model.Model) (Instance, error) {
	tokenCounts := make([]int8, model.PlaceCount)
	for i := range tokenCounts {
		tokenCounts[i] = 0
	}
	return Instance{
		Id:          createInstanceId(),
		TokenCounts: tokenCounts,
	}, nil
}

func (service *InstanceService) SaveInstance(instance Instance) error {
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
