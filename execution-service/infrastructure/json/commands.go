package json

import (
	"encoding/hex"
	"execution-service/domain"
	"execution-service/execution"
	"execution-service/proof"
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

type ProveInstantiationCommand struct {
	Model    Model  `json:"model"`
	Instance string `json:"instance"`
}

type ProveTransitionCommand struct {
	Model           Model  `json:"model"`
	CurrentInstance string `json:"currentInstance"`
	NextInstance    string `json:"nextInstance"`
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

func (cmd *ProveInstantiationCommand) ToProofCommand() (proof.ProveInstantiationCommand, error) {
	model, err := cmd.Model.ToDomainModel()
	if err != nil {
		return proof.ProveInstantiationCommand{}, err
	}
	return proof.ProveInstantiationCommand{
		Model:    model,
		Instance: cmd.Instance,
	}, nil
}

func (cmd *ProveTransitionCommand) ToProofCommand() (proof.ProveTransitionCommand, error) {
	model, err := cmd.Model.ToDomainModel()
	if err != nil {
		return proof.ProveTransitionCommand{}, err
	}
	return proof.ProveTransitionCommand{
		Model:           model,
		CurrentInstance: cmd.CurrentInstance,
		NextInstance:    cmd.NextInstance,
	}, nil
}

func (cmd *ProveTerminationCommand) ToProofCommand() (proof.ProveTerminationCommand, error) {
	model, err := cmd.Model.ToDomainModel()
	if err != nil {
		return proof.ProveTerminationCommand{}, err
	}
	return proof.ProveTerminationCommand{
		Model:    model,
		Instance: cmd.Instance,
	}, nil
}
