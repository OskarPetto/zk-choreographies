package execution

import (
	"execution-service/domain"
	"execution-service/proof"
)

type ExecutionResult struct {
	Instance domain.Instance
	Proof    proof.Proof
}
