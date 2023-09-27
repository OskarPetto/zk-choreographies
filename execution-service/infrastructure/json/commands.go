package json

import (
	"encoding/hex"
	"execution-service/domain"
	"execution-service/execution"
	"execution-service/proof"
)

type InstantiateModelCommand struct {
	Model      string   `json:"model"`
	PublicKeys []string `json:"publicKeys"`
}

type ExecuteTransitionCommand struct {
	Model      string `json:"model"`
	Instance   string `json:"instance"`
	Transition string `json:"transition"`
	Message    string `json:"message"`
}

type ProveInstantiationCommand struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
}

type ProveTransitionCommand struct {
	Model           string `json:"model"`
	CurrentInstance string `json:"currentInstance"`
	NextInstance    string `json:"nextInstance"`
}

type ProveTerminationCommand struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
}

func (cmd *InstantiateModelCommand) ToExecutionCommand() (execution.InstantiateModelCommand, error) {
	publicKeys := make([]domain.PublicKey, len(cmd.PublicKeys))
	for i, publicKey := range cmd.PublicKeys {
		bytes, err := hex.DecodeString(publicKey)
		if err != nil {
			return execution.InstantiateModelCommand{}, err
		}
		publicKeys[i] = domain.PublicKey{
			Value: bytes,
		}
	}
	return execution.InstantiateModelCommand{
		Model:      cmd.Model,
		PublicKeys: publicKeys,
	}, nil
}

func (cmd *ExecuteTransitionCommand) ToExecutionCommand() (execution.ExecuteTransitionCommand, error) {
	message, err := hex.DecodeString(cmd.Message)
	if err != nil {
		return execution.ExecuteTransitionCommand{}, err
	}
	return execution.ExecuteTransitionCommand{
		Model:      cmd.Model,
		Instance:   cmd.Instance,
		Transition: cmd.Transition,
		Message:    message,
	}, nil
}

func (cmd *ProveInstantiationCommand) ToProofCommand() (proof.ProveInstantiationCommand, error) {
	return proof.ProveInstantiationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
	}, nil
}

func (cmd *ProveTransitionCommand) ToProofCommand() (proof.ProveTransitionCommand, error) {
	return proof.ProveTransitionCommand{
		Model:           cmd.Model,
		CurrentInstance: cmd.CurrentInstance,
		NextInstance:    cmd.NextInstance,
	}, nil
}

func (cmd *ProveTerminationCommand) ToProofCommand() (proof.ProveTerminationCommand, error) {
	return proof.ProveTerminationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
	}, nil
}
