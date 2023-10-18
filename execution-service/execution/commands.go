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

type TerminateInstanceCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
	Identity domain.IdentityId
}

type TerminatedInstanceEvent struct {
	Proof prover.Proof
}

type SendMessageCommand struct {
	Model          domain.ModelId
	Instance       domain.InstanceId
	Transition     domain.TransitionId
	Identity       domain.IdentityId
	BytesMessage   []byte
	IntegerMessage *domain.IntegerType
}

type SentMessageEvent struct {
	Model            domain.ModelId
	CurrentInstance  domain.InstanceId
	Transition       domain.TransitionId
	NextInstance     domain.Instance
	SenderSignature  domain.Signature
	EncryptedMessage domain.Ciphertext
}

type ReceiveMessageCommand struct {
	Model            domain.ModelId
	CurrentInstance  domain.InstanceId
	Transition       domain.TransitionId
	Identity         domain.IdentityId
	NextInstance     domain.Instance
	SenderSignature  domain.Signature
	EncryptedMessage domain.Ciphertext
}

type ReceivedMessageEvent struct {
	Proof prover.Proof
}
