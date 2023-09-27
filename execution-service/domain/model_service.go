package domain

import "fmt"

type ModelService struct {
	isLoaded bool
	models   map[string]Model
}

var modelService ModelService

func NewModelService() ModelService {
	if !modelService.isLoaded {
		modelService = ModelService{
			isLoaded: true,
			models:   make(map[string]Model),
		}
	}
	return modelService
}

func (service *ModelService) SaveModel(model Model) {
	service.models[model.Id] = model
}

func (service *ModelService) FindModelById(id ModelId) (Model, error) {
	model, exists := service.models[id]
	if !exists {
		return Model{}, fmt.Errorf("model %s not found", id)
	}
	return model, nil
}
