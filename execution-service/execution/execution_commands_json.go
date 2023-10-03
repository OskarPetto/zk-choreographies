package execution

import (
	"execution-service/domain"
	"execution-service/message"
	"execution-service/utils"
)

type InstantiateModelCommandJson struct {
	Model      string   `json:"model"`
	PublicKeys []string `json:"publicKeys"`
}

type ExecuteTransitionCommandJson struct {
	Model                string                           `json:"model"`
	Instance             string                           `json:"instance"`
	Transition           string                           `json:"transition"`
	CreateMessageCommand message.CreateMessageCommandJson `json:"createMessageCommand"`
}

func (cmd *InstantiateModelCommandJson) ToExecutionCommand() (InstantiateModelCommand, error) {
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
	}, nil
}

func (cmd *ExecuteTransitionCommandJson) ToExecutionCommand() (ExecuteTransitionCommand, error) {
	return ExecuteTransitionCommand{
		Model:                cmd.Model,
		Instance:             cmd.Instance,
		Transition:           cmd.Transition,
		CreateMessageCommand: cmd.CreateMessageCommand.ToMessageCommand(),
	}, nil
}
