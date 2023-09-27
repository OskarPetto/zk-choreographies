package execution

import "execution-service/domain"

type InstantiateModelCommand struct {
	Model      domain.Model
	PublicKeys []domain.PublicKey
}

type ExecuteTransitionCommand struct {
	Model      domain.Model
	Instance   domain.InstanceId
	Transition domain.TransitionId
	Message    []byte
}

type ProveTerminationCommand struct {
	Model    domain.Model
	Instance domain.InstanceId
}
