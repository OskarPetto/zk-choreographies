package execution

import (
	"execution-service/instance"
	"execution-service/model"
)

type ExecutionService struct {
}

func (service *ExecutionService) InstantiateModel(model model.Model) instance.Instance {
	tokenCounts := make([]int8, model.PlaceCount)
	for i := range tokenCounts {
		tokenCounts[i] = 0
	}
	return instance.Instance{
		Model:       model.Id,
		TokenCounts: tokenCounts,
	}
}
