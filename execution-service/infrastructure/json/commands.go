package json

import (
	"encoding/hex"
	"execution-service/domain"
	"execution-service/execution"
)

type InstantiateModelCommand struct {
	Model      Model    `json:"model"`
	PublicKeys []string `json:"publicKeys"`
}

type ExecuteTransitionCommand struct {
	Model      Model  `json:"model"`
	Instance   string `json:"instance"`
	Transition string `json:"transition"`
	Message    string `json:"message"`
}

type ProveTerminationCommand struct {
	Model    Model  `json:"model"`
	Instance string `json:"instance"`
}

func (cmd *InstantiateModelCommand) ToExecutionCommand() (execution.InstantiateModelCommand, error) {
	model, err := cmd.Model.ToDomainModel()
	if err != nil {
		return execution.InstantiateModelCommand{}, err
	}
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
		Model:      model,
		PublicKeys: publicKeys,
	}, nil
}

func (cmd *ExecuteTransitionCommand) ToExecutionCommand() (execution.ExecuteTransitionCommand, error) {
	message, err := hex.DecodeString(cmd.Message)
	if err != nil {
		return execution.ExecuteTransitionCommand{}, err
	}
	model, err := cmd.Model.ToDomainModel()
	if err != nil {
		return execution.ExecuteTransitionCommand{}, err
	}
	return execution.ExecuteTransitionCommand{
		Model:      model,
		Instance:   cmd.Instance,
		Transition: cmd.Transition,
		Message:    message,
	}, nil
}

func (cmd *ProveTerminationCommand) ToExecutionCommand() (execution.ProveTerminationCommand, error) {
	model, err := cmd.Model.ToDomainModel()
	if err != nil {
		return execution.ProveTerminationCommand{}, err
	}
	return execution.ProveTerminationCommand{
		Model:    model,
		Instance: cmd.Instance,
	}, nil
}
