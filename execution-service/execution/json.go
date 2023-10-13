package execution

import (
	"execution-service/domain"
	"execution-service/prover"
	"execution-service/utils"
	"fmt"
)

type instantiateModelCommandJson struct {
	Model      string   `json:"model"`
	PublicKeys []string `json:"publicKeys"`
	Identity   uint     `json:"identity"`
}

type executeTransitionCommandJson struct {
	Model                string                    `json:"model"`
	Instance             string                    `json:"instance"`
	Transition           string                    `json:"transition"`
	Identity             uint                      `json:"identity"`
	CreateMessageCommand *createMessageCommandJson `json:"createMessageCommand,omitempty"`
}

type terminateInstanceCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	Identity uint   `json:"identity"`
}

type createMessageCommandJson struct {
	IntegerMessage *uint  `json:"integerMessage,omitempty"`
	BytesMessage   string `json:"bytesMessage,omitempty"`
}

type executionResultJson struct {
	Proof          prover.ProofJson `json:"proof"`
	EncryptedState string           `json:"encryptedState,omitempty"`
	PlainState     string           `json:"plainState,omitempty"`
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
	var createMessageCommand *domain.CreateMessageCommand = nil
	if cmd.CreateMessageCommand != nil {
		tmp, err := cmd.CreateMessageCommand.toDomainCommand(cmd.Model)
		if err != nil {
			return ExecuteTransitionCommand{}, nil
		}
		createMessageCommand = &tmp
	}
	return ExecuteTransitionCommand{
		Model:                cmd.Model,
		Instance:             cmd.Instance,
		Transition:           cmd.Transition,
		Identity:             cmd.Identity,
		CreateMessageCommand: createMessageCommand,
	}, nil
}

func (cmd *terminateInstanceCommandJson) ToExecutionCommand() (TerminateInstanceCommand, error) {
	return TerminateInstanceCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
		Identity: cmd.Identity,
	}, nil
}

func (cmd *createMessageCommandJson) toDomainCommand(modelId domain.ModelId) (domain.CreateMessageCommand, error) {
	model, err := utils.StringToBytes(modelId)
	if err != nil {
		return domain.CreateMessageCommand{}, err
	}
	modelFixed := [domain.HashSize]byte(model)

	bytesMessage, err := utils.StringToBytes(cmd.BytesMessage)
	if err == nil {
		return domain.CreateMessageCommand{
			Model: domain.Hash{
				Value: modelFixed,
			},
			IntegerMessage: nil,
			BytesMessage:   bytesMessage,
		}, nil
	} else if cmd.IntegerMessage != nil {
		intValue := domain.IntegerType(*cmd.IntegerMessage)
		return domain.CreateMessageCommand{
			Model: domain.Hash{
				Value: modelFixed,
			},
			IntegerMessage: &intValue,
			BytesMessage:   bytesMessage,
		}, nil
	}
	return domain.CreateMessageCommand{}, fmt.Errorf("createMessageCommand of model %s could not be parsed", modelId)
}

func ToJson(result ExecutionResult) executionResultJson {
	encryptedState := ""
	if result.EncryptedState != nil {
		encryptedState = utils.BytesToString(result.EncryptedState.Value)
	}
	plainState := ""
	if result.PlainState != nil {
		plainState = string(result.PlainState.Value)
	}
	return executionResultJson{
		Proof:          result.Proof.ToJson(),
		EncryptedState: encryptedState,
		PlainState:     plainState,
	}
}
