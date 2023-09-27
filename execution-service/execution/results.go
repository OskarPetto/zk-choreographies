package execution

import (
	"execution-service/domain"
	"execution-service/proof"
)

type InstantiateModelResult struct {
	Instance domain.Instance
	Proof    proof.Proof
}

type ExecuteTransitionResult struct {
	Instance domain.Instance
	Proof    proof.Proof
}

type ProveTerminationResult struct {
	Proof proof.Proof
}
