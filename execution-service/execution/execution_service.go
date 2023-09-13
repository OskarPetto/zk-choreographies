package execution

import (
	inst "execution-service/instance"
	"execution-service/model"
	"fmt"
)

type ExecutionService struct {
}

func (service *ExecutionService) InstantiateModel(model model.Model) inst.Instance {
	tokenCounts := make([]int8, model.PlaceCount)
	for i := range tokenCounts {
		tokenCounts[i] = 0
	}
	tokenCounts[model.StartPlace] = 1
	return inst.Instance{
		Model:       model.Id,
		TokenCounts: tokenCounts,
	}
}

func (service *ExecutionService) ExecuteTransition(instance inst.Instance, transition model.Transition) (inst.Instance, error) {
	if !isTransitionExecutable(instance, transition) {
		return inst.Instance{}, fmt.Errorf("transition %+v is not executable in instance %+v", transition, instance)
	}
	newInstance := copyInstance(instance)
	for _, placeId := range transition.IncomingPlaces {
		newInstance.TokenCounts[placeId] -= 1
	}
	for _, placeId := range transition.OutgoingPlaces {
		newInstance.TokenCounts[placeId] += 1
	}
	return newInstance, nil
}

func isTransitionExecutable(instance inst.Instance, transition model.Transition) bool {
	for _, placeId := range transition.IncomingPlaces {
		if instance.TokenCounts[placeId] < 1 {
			return false
		}
	}
	return true
}

func copyInstance(instance inst.Instance) inst.Instance {
	newInstance := inst.Instance{
		Id:          instance.Id,
		Model:       instance.Model,
		TokenCounts: make([]int8, len(instance.TokenCounts)),
	}
	copy(newInstance.TokenCounts, instance.TokenCounts)
	return newInstance
}
