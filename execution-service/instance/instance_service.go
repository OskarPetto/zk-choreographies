package instance

import (
	"execution-service/domain"
	"fmt"
	"sort"
)

type InstanceService struct {
	instances map[string]domain.Instance
}

func NewInstanceService() InstanceService {
	return InstanceService{
		instances: make(map[string]domain.Instance),
	}
}

func (service *InstanceService) SaveInstance(instance domain.Instance) {
	service.instances[instance.Id()] = instance
}

func (service *InstanceService) DeleteInstance(id domain.InstanceId) {
	delete(service.instances, id)
}

func (service *InstanceService) FindInstanceById(id domain.InstanceId) (domain.Instance, error) {
	instance, exists := service.instances[id]
	if !exists {
		return domain.Instance{}, fmt.Errorf("instance %s not found", id)
	}
	return instance, nil
}

func (service *InstanceService) FindInstancesByModel(model domain.ModelId) []domain.Instance {
	instances := make([]domain.Instance, 0, len(service.instances))
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
