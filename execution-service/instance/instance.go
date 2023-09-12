package instance

import (
	"execution-service/model"
)

type InstanceId string

type Instance struct {
	Id          InstanceId
	Model       model.ModelId
	TokenCounts []int8
}
