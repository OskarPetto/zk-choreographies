package execution

import (
	"execution-service/domain"
	"execution-service/prover"
)

type ExecutionResult struct {
	Proof    prover.Proof
	Instance domain.Instance
}
