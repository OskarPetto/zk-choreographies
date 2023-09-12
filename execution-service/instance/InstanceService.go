package instance

import (
	"execution-service/model"

	"github.com/google/uuid"
)

type InstanceService interface {
	InstantiateModel(model model.Model) (Instance, error)
}

func InstantiateModel(model model.Model) (Instance, error) {
	var tokenCounts []int8 = make([]int8, model.PlaceCount)
	for i := range tokenCounts {
		tokenCounts[i] = 0
	}
	return Instance{
		Id:          createInstanceId(),
		TokenCounts: tokenCounts,
	}, nil
}

func createInstanceId() InstanceId {
	return InstanceId(uuid.New().String())
}
