package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/prover"
	"execution-service/utils"
)

type instantiateModelCommandJson struct {
	Model      string   `json:"model"`
	PublicKeys []string `json:"publicKeys"`
	Identity   uint     `json:"identity"`
}

type executeTransitionCommandJson struct {
	Model      string `json:"model"`
	Instance   string `json:"instance"`
	Transition string `json:"transition"`
	Identity   uint   `json:"identity"`
}

type terminateInstanceCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	Identity uint   `json:"identity"`
}

type executionResultJson struct {
	Proof    prover.ProofJson      `json:"proof"`
	Instance instance.InstanceJson `json:"instance"`
}

func (cmd *instantiateModelCommandJson) ToExecutionCommand() (InstantiateModelCommand, error) {
	publicKeys := make([]domain.PublicKey, len(cmd.PublicKeys))
	for i, publicKey := range cmd.PublicKeys {
		bytes, err := utils.StringToBytes(publicKey)
		if err != nil {
			return InstantiateModelCommand{}, err
		}
		publicKeys[i] = domain.PublicKey{
			Value: bytes,
		}
	}
	return InstantiateModelCommand{
		Model:      cmd.Model,
		PublicKeys: publicKeys,
		Identity:   cmd.Identity,
	}, nil
}

func (cmd *executeTransitionCommandJson) ToExecutionCommand() (ExecuteTransitionCommand, error) {
	return ExecuteTransitionCommand{
		Model:      cmd.Model,
		Instance:   cmd.Instance,
		Transition: cmd.Transition,
		Identity:   cmd.Identity,
	}, nil
}

func (cmd *terminateInstanceCommandJson) ToExecutionCommand() (TerminateInstanceCommand, error) {
	return TerminateInstanceCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
		Identity: cmd.Identity,
	}, nil
}

func ToJson(result ExecutionResult) executionResultJson {
	return executionResultJson{
		Proof:    result.Proof.ToJson(),
		Instance: instance.ToJson(result.Instance),
	}
}
