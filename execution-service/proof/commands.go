package proof

import (
	"execution-service/domain"
)

type ProveInstantiationCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
}

type ProveTransitionCommand struct {
	Model           domain.ModelId
	CurrentInstance domain.InstanceId
	NextInstance    domain.InstanceId
}

type ProveTerminationCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
}
