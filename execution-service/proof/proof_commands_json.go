package proof

type ProveInstantiationCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
}

type ProveTransitionCommandJson struct {
	Model           string `json:"model"`
	CurrentInstance string `json:"currentInstance"`
	NextInstance    string `json:"nextInstance"`
}

type ProveTerminationCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
}

func (cmd *ProveInstantiationCommandJson) ToProofCommand() (ProveInstantiationCommand, error) {
	return ProveInstantiationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
	}, nil
}

func (cmd *ProveTransitionCommandJson) ToProofCommand() (ProveTransitionCommand, error) {
	return ProveTransitionCommand{
		Model:           cmd.Model,
		CurrentInstance: cmd.CurrentInstance,
		NextInstance:    cmd.NextInstance,
	}, nil
}

func (cmd *ProveTerminationCommandJson) ToProofCommand() (ProveTerminationCommand, error) {
	return ProveTerminationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
	}, nil
}
