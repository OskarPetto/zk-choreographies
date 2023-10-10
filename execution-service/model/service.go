package model

import (
	"execution-service/domain"
	"fmt"
	"sort"
)

type ModelService struct {
	models map[domain.ModelId]domain.Model
}

func NewModelService() ModelService {
	return ModelService{
		models: make(map[string]domain.Model),
	}
}

func (service *ModelService) CreateModel(model domain.Model) domain.Model {
	model.ComputeHash()
	service.models[model.Id()] = model
	return model
}

func (service *ModelService) ImportModel(model domain.Model) {
	service.models[model.Id()] = model
}

func (service *ModelService) FindModelById(modelId domain.ModelId) (domain.Model, error) {
	model, exists := service.models[modelId]
	if exists {
		return model, nil
	}
	return domain.Model{}, fmt.Errorf("model %s not found", modelId)
}

func (service *ModelService) FindModelsByChoreography(choreography string) []domain.Model {
	models := make([]domain.Model, 0, len(service.models))
	for _, model := range service.models {
		if model.Choreography == choreography {
			models = append(models, model)
		}
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].CreatedAt > models[j].CreatedAt
	})
	return models
}
