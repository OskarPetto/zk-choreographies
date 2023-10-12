package execution

import (
	"execution-service/domain"
)

type CreateMessageCommand struct {
	BytesMessage   []byte
	IntegerMessage domain.IntegerType
}

type InstantiateModelCommand struct {
	Model      domain.ModelId
	PublicKeys []domain.PublicKey
	Identity   domain.IdentityId
}

type ExecuteTransitionCommand struct {
	Model                domain.ModelId
	Instance             domain.InstanceId
	Transition           domain.TransitionId
	Identity             domain.IdentityId
	CreateMessageCommand *CreateMessageCommand
}

type TerminateInstanceCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
	Identity domain.IdentityId
}
