package model

import (
	"bytes"
	"execution-service/domain"
	"execution-service/instance"
	"fmt"
	"sort"
)

type ModelService struct {
	models          map[domain.ModelId]domain.Model
	InstanceService instance.InstanceService
}

func NewModelService(instanceService instance.InstanceService) ModelService {
	return ModelService{
		models:          make(map[string]domain.Model),
		InstanceService: instanceService,
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

func (service *ModelService) ImportModel(cmd ImportModelCommand) error {
	if !cmd.Model.HasValidHash() {
		return fmt.Errorf("model %s has invalid hash", cmd.Model.Id())
	}
	if !bytes.Equal(cmd.Instance.Model.Value[:], cmd.Model.Hash.Hash.Value[:]) {
		return fmt.Errorf("model %s does not fit instance %s", cmd.Model.Id(), cmd.Instance.Id())
	}
	err := service.InstanceService.ImportInstance(cmd.Instance)
	if err != nil {
		return err
	}
	service.models[cmd.Model.Id()] = cmd.Model
	return nil
}

func (service *ModelService) CreateModel(model domain.Model) domain.SaltedHash {
	model.UpdateHash()
	service.models[model.Id()] = model
	return model.Hash
}

func (service *ModelService) DeleteModel(model domain.Model) {
	delete(service.models, model.Id())
}
