package execution

import (
	"execution-service/authentication"
	"execution-service/domain"
)

type ExecutionService struct {
	isLoaded         bool
	signatureService authentication.SignatureService
	instanceService  domain.InstanceService
	hashService      domain.HashService
	modelService     domain.ModelService
}

var executionService ExecutionService

func NewExecutionService() ExecutionService {
	if !executionService.isLoaded {
		executionService = ExecutionService{
			isLoaded:        true,
			instanceService: domain.NewInstanceService(),
			hashService:     domain.NewHashService(),
			modelService:    domain.ModelServiceImpl,
		}
	}
	return executionService
}

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (domain.Instance, error) {
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return domain.Instance{}, err
	}
	modelHash := domain.HashModel(model)
	instanceResult, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return domain.Instance{}, err
	}
	service.hashService.SaveModelHash(model.Id, modelHash)
	service.instanceService.SaveInstance(instanceResult)
	return instanceResult, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (domain.Instance, error) {
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return domain.Instance{}, err
	}
	_, err = service.hashService.FindHashByModelId(model.Id)
	if err != nil {
		return domain.Instance{}, err
	}
	currentInstance, err := service.instanceService.FindInstanceById(cmd.Instance)
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
	return nextInstance, nil
}
