package execution

import (
	"execution-service/domain"
	"execution-service/prover"
)

type ExecutionResult struct {
	Instance       domain.Instance
	Proof          prover.Proof
	EncryptedState domain.Chiphertext
}
