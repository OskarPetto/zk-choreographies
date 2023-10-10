package instance

import (
	"execution-service/domain"
	"execution-service/message"
	"execution-service/model"
	"fmt"
	"sort"
)

type InstanceService struct {
	instances      map[string]domain.Instance
	ModelService   model.ModelService
	MessageService message.MessageService
}

func InitializeInstanceService() InstanceService {
	modelService := model.NewModelService()
	messageService := message.NewMessageService()
	return NewInstanceService(modelService, messageService)
}

func NewInstanceService(modelService model.ModelService, messageService message.MessageService) InstanceService {
	return InstanceService{
		instances:      make(map[string]domain.Instance),
		ModelService:   modelService,
		MessageService: messageService,
	}
}

func (service *InstanceService) ImportInstance(instance domain.Instance) {
	service.instances[instance.Id()] = instance
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
		return instances[i].CreatedAt > instances[j].CreatedAt
	})
	return instances
}

func (service *InstanceService) InstantiateModel(cmd InstantiateModelCommand) (domain.Instance, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return domain.Instance{}, err
	}
	instanceResult, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return domain.Instance{}, err
	}
	service.ImportInstance(instanceResult)
	return instanceResult, nil
}

func (service *InstanceService) ExecuteTransition(cmd ExecuteTransitionCommand) (domain.Instance, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return domain.Instance{}, err
	}
	currentInstance, err := service.FindInstanceById(cmd.Instance)
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
	service.ImportInstance(nextInstance)
	return nextInstance, nil
}
