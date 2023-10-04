package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
)

type ExecutionService struct {
	InstanceService instance.InstanceService
	ModelService    model.ModelService
	MessageService  message.MessageService
}

func InitializeExecutionService() ExecutionService {
	instanceService := instance.NewInstanceService()
	modelService := model.NewModelService()
	messageService := message.NewMessageService()
	return NewExecutionService(instanceService, modelService, messageService)
}

func NewExecutionService(instanceService instance.InstanceService, modelService model.ModelService, messageService message.MessageService) ExecutionService {
	return ExecutionService{
		InstanceService: instanceService,
		ModelService:    modelService,
		MessageService:  messageService,
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
	service.InstanceService.ImportInstance(instanceResult)
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
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return domain.Instance{}, err
	}
	var nextInstance domain.Instance
	if cmd.CreateMessageCommand != nil {
		message := service.MessageService.CreateMessage(*cmd.CreateMessageCommand)
		nextInstance, err = currentInstance.ExecuteTransition(transition, constraintInput, message.Hash)
	} else {
		nextInstance, err = currentInstance.ExecuteTransition(transition, constraintInput, domain.EmptyHash())
	}
	if err != nil {
		return domain.Instance{}, err
	}
	service.InstanceService.ImportInstance(nextInstance)
	return nextInstance, nil
}
