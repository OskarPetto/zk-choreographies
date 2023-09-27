package execution

import (
	"execution-service/domain"
	"execution-service/proof"
)

type InstantiationResult struct {
	Instance domain.Instance
	Proof    proof.Proof
}

type TransitionResult struct {
	Instance domain.Instance
	Proof    proof.Proof
}

type TerminationResult struct {
	Proof proof.Proof
}
