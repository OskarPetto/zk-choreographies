//go:build wireinject
// +build wireinject

package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/model"

	"github.com/google/wire"
)

type ExecutionService struct {
	InstanceService instance.InstanceService
	ModelService    model.ModelService
}

func InitializeExecutionService(modelPort model.ModelPort) ExecutionService {
	wire.Build(instance.NewInstanceService, model.NewModelService, NewExecutionService)
	return ExecutionService{}
}

func NewExecutionService(instanceService instance.InstanceService, modelService model.ModelService) ExecutionService {
	return ExecutionService{
		InstanceService: instanceService,
		ModelService:    modelService,
	}
}

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (domain.Instance, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return domain.Instance{}, err
	}
	instanceResult, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return domain.Instance{}, err
	}
	service.InstanceService.SaveInstance(instanceResult)
	return instanceResult, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (domain.Instance, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return domain.Instance{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return domain.Instance{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return domain.Instance{}, err
	}
	var nextInstance domain.Instance
	if len(cmd.Message) == 0 {
		nextInstance, err = currentInstance.ExecuteTransition(transition)
	} else {
		nextInstance, err = currentInstance.ExecuteTransitionWithMessage(transition, cmd.Message)
	}
	if err != nil {
		return domain.Instance{}, err
	}
	service.InstanceService.SaveInstance(nextInstance)
	return nextInstance, nil
}
