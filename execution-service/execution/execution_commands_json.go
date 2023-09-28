package execution

import (
	"execution-service/domain"
	"execution-service/utils"
)

type InstantiateModelCommandJson struct {
	PublicKeys []string `json:"publicKeys"`
}

type ExecuteTransitionCommandJson struct {
	Transition string `json:"transition"`
	Message    string `json:"message"`
}

func (cmd *InstantiateModelCommandJson) ToExecutionCommand(modelId domain.ModelId) (InstantiateModelCommand, error) {
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
		Model:      modelId,
		PublicKeys: publicKeys,
	}, nil
}

func (cmd *ExecuteTransitionCommandJson) ToExecutionCommand(modelId domain.ModelId, instanceId domain.InstanceId) (ExecuteTransitionCommand, error) {
	message, err := utils.StringToBytes(cmd.Message)
	if err != nil {
		return ExecuteTransitionCommand{}, err
	}
	return ExecuteTransitionCommand{
		Model:      modelId,
		Instance:   instanceId,
		Transition: cmd.Transition,
		Message:    message,
	}, nil
}
