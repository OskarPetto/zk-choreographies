package json

import (
	"execution-service/execution"
	"execution-service/proof"
)

type Proof struct {
	A           [2]string    `json:"a"`
	B           [2][2]string `json:"b"`
	C           [2]string    `json:"c"`
	PublicInput []string     `json:"publicInput"`
}

type ExecutionResult struct {
	Instance Instance
	Proof    Proof
}

type ProofResult struct {
	Proof Proof
}

func FromExecutionResult(executionResult execution.ExecutionResult) ExecutionResult {
	return ExecutionResult{
		Instance: FromDomainInstance(executionResult.Instance),
		Proof:    fromDomainProof(executionResult.Proof),
	}
}

func fromDomainProof(proof proof.Proof) Proof {
	publicInputs := make([]string, len(proof.PublicInput))
	for i, publicInput := range proof.PublicInput {
		publicInputs[i] = publicInput.String()
	}
	return Proof{
		A: [2]string{
			proof.A[0].String(),
			proof.A[1].String(),
		},
		B: [2][2]string{
			{
				proof.B[0][0].String(),
				proof.B[0][1].String(),
			},
			{
				proof.B[1][0].String(),
				proof.B[1][1].String(),
			},
		},
		C: [2]string{
			proof.C[0].String(),
			proof.C[1].String(),
		},
		PublicInput: publicInputs,
	}
}
