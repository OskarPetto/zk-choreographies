package prover

import (
	"execution-service/domain"
)

type ProveInstantiationCommand struct {
	Model     domain.Model
	Instance  domain.Instance
	Signature domain.Signature
}

type ProveTransitionCommand struct {
	Model                          domain.Model
	CurrentInstance                domain.Instance
	NextInstance                   domain.Instance
	Transition                     domain.Transition
	InitiatingParticipantSignature domain.Signature
	RespondingParticipantSignature *domain.Signature
	ConstraintInput                domain.ConstraintInput
}

type ProveTerminationCommand struct {
	Model     domain.Model
	Instance  domain.Instance
	Signature domain.Signature
}
