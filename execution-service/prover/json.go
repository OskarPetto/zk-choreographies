package prover

import "execution-service/domain"

type ProveInstantiationCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	Identity uint   `json:"identity"`
}

type ProveTransitionCommandJson struct {
	Model           string `json:"model"`
	CurrentInstance string `json:"currentInstance"`
	NextInstance    string `json:"nextInstance"`
	Transition      string `json:"transition"`
	Identity        uint   `json:"identity"`
}

type ProveTerminationCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	EndPlace uint   `json:"endPlace"`
	Identity uint   `json:"identity"`
}

type ProofJson struct {
	Value [8]string `json:"value"`
	Input []string  `json:"input"`
}

func (cmd *ProveInstantiationCommandJson) ToProofCommand() (ProveInstantiationCommand, error) {
	return ProveInstantiationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
		Identity: cmd.Identity,
	}, nil
}

func (cmd *ProveTransitionCommandJson) ToProofCommand() (ProveTransitionCommand, error) {
	return ProveTransitionCommand{
		Model:           cmd.Model,
		CurrentInstance: cmd.CurrentInstance,
		NextInstance:    cmd.NextInstance,
		Transition:      cmd.Transition,
		Identity:        cmd.Identity,
	}, nil
}

func (cmd *ProveTerminationCommandJson) ToProofCommand() (ProveTerminationCommand, error) {
	return ProveTerminationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
		EndPlace: domain.PlaceId(cmd.EndPlace),
		Identity: cmd.Identity,
	}, nil
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
