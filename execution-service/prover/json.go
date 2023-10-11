package prover

type ProofJson struct {
	Value [8]string `json:"value"`
	Input []string  `json:"input"`
}

func (proof Proof) ToJson() ProofJson {
	publicInputs := make([]string, len(proof.Input))
	for i, publicInput := range proof.Input {
		publicInputs[i] = publicInput.String()
	}
	return ProofJson{
		Value: [8]string{
			proof.Value[0].String(),
			proof.Value[1].String(),
			proof.Value[2].String(),
			proof.Value[3].String(),
			proof.Value[4].String(),
			proof.Value[5].String(),
			proof.Value[6].String(),
			proof.Value[7].String(),
		},
		Input: publicInputs,
	}
}
