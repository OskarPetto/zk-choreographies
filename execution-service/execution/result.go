package execution

import (
	"execution-service/domain"
	"execution-service/prover"
)

type ExecutionResult struct {
	Proof          prover.Proof
	EncryptedState *domain.Ciphertext
	PlainState     *domain.Plaintext
}
