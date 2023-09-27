package execution

import "execution-service/domain"

type InstantiateModelCommand struct {
	Model      domain.ModelId
	PublicKeys []domain.PublicKey
}

type ExecuteTransitionCommand struct {
	Model      domain.ModelId
	Instance   domain.InstanceId
	Transition domain.TransitionId
	Message    []byte
}
