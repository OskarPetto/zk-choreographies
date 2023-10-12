package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/prover"
	"execution-service/utils"
)

type instantiateModelCommandJson struct {
	Model      model.ModelJson `json:"model"`
	PublicKeys []string        `json:"publicKeys"`
	Identity   uint            `json:"identity"`
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
	IntegerMessage uint   `json:"integerMessage,omitempty"`
	BytesMessage   string `json:"bytesMessage,omitempty"`
}

type executionResultJson struct {
	Instance       instance.InstanceJson `json:"instance"`
	Proof          prover.ProofJson      `json:"proof"`
	EncryptedState string                `json:"encryptedState"`
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
	model, err := cmd.Model.ToModel()
	if err != nil {
		return InstantiateModelCommand{}, err
	}
	return InstantiateModelCommand{
		Model:      model,
		PublicKeys: publicKeys,
		Identity:   cmd.Identity,
	}, nil
}

func (cmd *executeTransitionCommandJson) ToExecutionCommand() (ExecuteTransitionCommand, error) {
	var createMessageCommand *CreateMessageCommand = nil
	if cmd.CreateMessageCommand != nil {
		tmp := cmd.CreateMessageCommand.ToExecutionCommand()
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

func (cmd *createMessageCommandJson) ToExecutionCommand() CreateMessageCommand {
	bytesMessage, err := utils.StringToBytes(cmd.BytesMessage)
	if err != nil {
		return CreateMessageCommand{
			IntegerMessage: int32(cmd.IntegerMessage),
			BytesMessage:   []byte{},
		}
	}

	return CreateMessageCommand{
		IntegerMessage: domain.EmptyIntegerMessage,
		BytesMessage:   bytesMessage,
	}
}

func ToJson(result ExecutionResult) executionResultJson {
	return executionResultJson{
		Instance:       instance.ToJson(result.Instance),
		Proof:          result.Proof.ToJson(),
		EncryptedState: utils.BytesToString(result.EncryptedState.Value),
	}
}
