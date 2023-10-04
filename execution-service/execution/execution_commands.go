package execution

import (
	"execution-service/domain"
	"execution-service/message"
)

type InstantiateModelCommand struct {
	Model      domain.ModelId
	PublicKeys []domain.PublicKey
}

type ExecuteTransitionCommand struct {
	Model                domain.ModelId
	Instance             domain.InstanceId
	Transition           domain.TransitionId
	CreateMessageCommand *message.CreateMessageCommand
}
