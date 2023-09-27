package domain

import (
	"fmt"
	"sort"
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
	service.instances[instance.Id()] = instance
}

func (service *InstanceService) DeleteInstance(id InstanceId) {
	delete(service.instances, id)
}

func (service *InstanceService) FindInstanceById(id InstanceId) (Instance, error) {
	instance, exists := service.instances[id]
	if !exists {
		return Instance{}, fmt.Errorf("instance %s not found", id)
	}
	return instance, nil
}

func (service *InstanceService) FindInstancesByModel(model ModelId) []Instance {
	instances := make([]Instance, 0, len(service.instances))
	for _, instance := range service.instances {
		if instance.Model == model {
			instances = append(instances, instance)
		}
	}
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].UpdatedAt > instances[j].UpdatedAt
	})
	return instances
}
