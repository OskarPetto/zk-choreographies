package proof

import (
	"execution-service/authentication"
	"execution-service/domain"
)

type ProveInstantiationCommand struct {
	ModelHash domain.Hash
	Model     domain.Model
	Instance  domain.Instance
	Signature authentication.Signature
}

type ProveTransitionCommand struct {
	ModelHash       domain.Hash
	Model           domain.Model
	CurrentInstance domain.Instance
	NextInstance    domain.Instance
	NextSignature   authentication.Signature
}

type ProveTerminationCommand struct {
	ModelHash domain.Hash
	Model     domain.Model
	Instance  domain.Instance
	Signature authentication.Signature
}
