package proof

import (
	"execution-service/domain"
)

type ProveInstantiationCommand struct {
	Model    domain.Model
	Instance domain.InstanceId
}

type ProveTransitionCommand struct {
	Model           domain.Model
	CurrentInstance domain.InstanceId
	NextInstance    domain.InstanceId
}

type ProveTerminationCommand struct {
	Model    domain.Model
	Instance domain.InstanceId
}
