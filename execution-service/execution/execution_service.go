package execution

import (
	"proof-service/proof"
)

type ExecutionService struct {
	isLoaded     bool
	proofService proof.ProofService
}

var executionService ExecutionService

func NewExecutionService() ExecutionService {
	if !executionService.isLoaded {
		executionService = ExecutionService{
			isLoaded:     true,
			proofService: proof.NewProofService(),
		}
	}
	return executionService
}
