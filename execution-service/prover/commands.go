package prover

import (
	"execution-service/domain"
)

type ProveInstantiationCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
	Identity domain.IdentityId
}

type ProveTransitionCommand struct {
	Model           domain.ModelId
	CurrentInstance domain.InstanceId
	NextInstance    domain.InstanceId
	Transition      domain.TransitionId
	Identity        domain.IdentityId
}

type ProveTerminationCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
	EndPlace domain.PlaceId
	Identity domain.IdentityId
}
