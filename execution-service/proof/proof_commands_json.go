package proof

type ProveInstantiationCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	Identity uint   `json:"identity"`
}

type ProveTransitionCommandJson struct {
	Model           string `json:"model"`
	CurrentInstance string `json:"currentInstance"`
	NextInstance    string `json:"nextInstance"`
	Identity        uint   `json:"identity"`
}

type ProveTerminationCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	Identity uint   `json:"identity"`
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
		Identity:        cmd.Identity,
	}, nil
}

func (cmd *ProveTerminationCommandJson) ToProofCommand() (ProveTerminationCommand, error) {
	return ProveTerminationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
		Identity: cmd.Identity,
	}, nil
}
