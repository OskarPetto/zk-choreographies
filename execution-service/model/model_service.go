package model

import (
	"execution-service/domain"
	"fmt"
)

type ModelService struct {
	models    map[domain.ModelId]domain.Model
	modelPort ModelPort
}

func NewModelService(modelPort ModelPort) ModelService {
	return ModelService{
		models:    make(map[string]domain.Model),
		modelPort: modelPort,
	}
}

func (service *ModelService) SaveModel(model domain.Model) {
	service.models[model.Id] = model
}

func (service *ModelService) FindModelById(modelId domain.ModelId) (domain.Model, error) {
	model, exists := service.models[modelId]
	if exists {
		return model, nil
	}
	model, err := service.modelPort.FindModelById(modelId)
	if err != nil {
		return domain.Model{}, fmt.Errorf("model %s not found in model-service", modelId)
	}
	model.ComputeHash()
	service.SaveModel(model)
	return model, nil
}
