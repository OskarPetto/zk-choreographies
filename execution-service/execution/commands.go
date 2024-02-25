package execution

import (
	"execution-service/domain"
	"execution-service/prover"
)

type InstantiateModelCommand struct {
	Model      domain.ModelId
	PublicKeys []domain.PublicKey
	Identity   domain.IdentityId
}

type InstantiatedModelEvent struct {
	Instance domain.Instance
	Proof    prover.Proof
}

type ExecuteTransitionCommand struct {
	Model      domain.ModelId
	Instance   domain.InstanceId
	Transition domain.TransitionId
	Identity   domain.IdentityId
}

type ExecutedTransitionEvent struct {
	Instance domain.Instance
	Proof    prover.Proof
}

type ProveTerminationCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
	Identity domain.IdentityId
}

type ProvedTerminationEvent struct {
	Proof prover.Proof
}

type CreateInitiatingMessageCommand struct {
	Model          domain.ModelId
	Instance       domain.InstanceId
	Transition     domain.TransitionId
	BytesMessage   []byte
	IntegerMessage *domain.IntegerType
}

type CreatedInitiatingMessageEvent struct {
	Model              domain.Model
	CurrentInstance    domain.Instance
	Transition         domain.TransitionId
	InintiatingMessage domain.Message
}

type ReceiveInitiatingMessageCommand struct {
	Model             domain.Model
	CurrentInstance   domain.Instance
	Transition        domain.TransitionId
	Identity          domain.IdentityId
	InitiatingMessage domain.Message
	BytesMessage      []byte
	IntegerMessage    *domain.IntegerType
}

type ReceivedInitiatingMessageEvent struct {
	Model                          domain.ModelId
	CurrentInstance                domain.InstanceId
	Transition                     domain.TransitionId
	InitiatingMessage              domain.MessageId
	NextInstance                   domain.Instance
	RespondingMessage              *domain.Message
	RespondingParticipantSignature domain.Signature
}

type ProveMessageExchangeCommand struct {
	Model                          domain.ModelId
	CurrentInstance                domain.InstanceId
	Transition                     domain.TransitionId
	Identity                       domain.IdentityId
	InitiatingMessage              domain.MessageId
	NextInstance                   domain.Instance
	RespondingMessage              *domain.Message
	RespondingParticipantSignature domain.Signature
}

type ProvedMessageExchangeEvent struct {
	Instance domain.Instance
	Proof    prover.Proof
}
