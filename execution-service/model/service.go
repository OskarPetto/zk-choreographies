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

func (service *ModelService) FindModelById(modelId domain.ModelId) (domain.Model, error) {
	model, exists := service.models[modelId]
	if exists {
		return model, nil
	}
	return domain.Model{}, fmt.Errorf("model %s not found", modelId)
}

func (service *ModelService) FindAllModels() []domain.Model {
	models := make([]domain.Model, 0, len(service.models))
	for _, model := range service.models {
		models = append(models, model)
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].CreatedAt > models[j].CreatedAt
	})
	return models
}

func (service *ModelService) ImportModel(model domain.Model) error {
	if !model.HasValidHash() {
		return fmt.Errorf("model %s has invalid hash", model.Id())
	}
	service.saveModel(model)
	return nil
}

func (service *ModelService) CreateModel(model domain.Model) domain.SaltedHash {
	model.UpdateHash()
	service.saveModel(model)
	return model.SaltedHash
}

func (service *ModelService) DeleteModel(model domain.Model) {
	delete(service.models, model.Id())
}

func (service *ModelService) saveModel(model domain.Model) {
	service.models[model.Id()] = model
}
