package instance

import (
	"bytes"
	"execution-service/domain"
	"execution-service/model"
	"execution-service/utils"
	"fmt"
	"sort"
)

type InstanceService struct {
	ModelService model.ModelService
	instances    map[string]domain.Instance
}

func NewInstanceService(modelService model.ModelService) InstanceService {
	return InstanceService{
		instances:    make(map[string]domain.Instance),
		ModelService: modelService,
	}
}

func (service *InstanceService) FindInstanceById(id domain.InstanceId) (domain.Instance, error) {
	instance, exists := service.instances[id]
	if !exists {
		return domain.Instance{}, fmt.Errorf("instance %s not found", id)
	}
	return instance, nil
}

func (service *InstanceService) FindInstancesByModel(model domain.ModelId) []domain.Instance {
	modelHash, err := utils.StringToBytes(model)
	if err != nil {
		return []domain.Instance{}
	}
	instances := make([]domain.Instance, 0, len(service.instances))
	for _, instance := range service.instances {
		if bytes.Equal(instance.Model.Value[:], modelHash) {
			instances = append(instances, instance)
		}
	}
	sort.Slice(instances, func(i, j int) bool {
		return instances[i].CreatedAt > instances[j].CreatedAt
	})
	return instances
}

func (service *InstanceService) ImportInstance(instance domain.Instance) error {
	if !instance.HasValidHash() {
		return fmt.Errorf("instance %s has invalid hash", instance.Id())
	}
	modelId := utils.BytesToString(instance.Model.Value[:])
	_, err := service.ModelService.FindModelById(modelId)
	if err != nil {
		return err
	}
	service.SaveInstance(instance)
	return nil
}

func (service *InstanceService) SaveInstance(instance domain.Instance) {
	service.instances[instance.Id()] = instance
}
