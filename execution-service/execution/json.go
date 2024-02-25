package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/prover"
	"execution-service/signature"
	"execution-service/utils"
)

type instantiateModelCommandJson struct {
	Model      string   `json:"model"`
	PublicKeys []string `json:"publicKeys"`
	Identity   uint     `json:"identity"`
}

type instantiatedModelEventJson struct {
	Instance instance.InstanceJson `json:"instance"`
	Proof    prover.ProofJson      `json:"proof"`
}

type executeTransitionCommandJson struct {
	Model      string `json:"model"`
	Instance   string `json:"instance"`
	Transition string `json:"transition"`
	Identity   uint   `json:"identity"`
}

type executedTransitionEventJson struct {
	Instance instance.InstanceJson `json:"instance"`
	Proof    prover.ProofJson      `json:"proof"`
}

type proveTerminationCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	Identity uint   `json:"identity"`
}

type provedTerminationEventJson struct {
	Proof prover.ProofJson `json:"proof"`
}

type createInitiatingMessageCommandJson struct {
	Model          string `json:"model"`
	Instance       string `json:"instance"`
	Transition     string `json:"transition"`
	BytesMessage   []byte `json:"bytesMessage,omitempty"`
	IntegerMessage *uint  `json:"integerMessage,omitempty"`
}

type createdInitiatingMessageEventJson struct {
	Model             model.ModelJson       `json:"model"`
	CurrentInstance   instance.InstanceJson `json:"currentInstance"`
	Transition        string                `json:"transition"`
	InitiatingMessage message.MessageJson   `json:"initiatingMessage"`
}

type receiveInitiatingMessageCommandJson struct {
	Model             model.ModelJson       `json:"model"`
	CurrentInstance   instance.InstanceJson `json:"currentInstance"`
	Transition        string                `json:"transition"`
	InitiatingMessage message.MessageJson   `json:"initiatingMessage"`
	Identity          uint                  `json:"identity"`
	BytesMessage      []byte                `json:"bytesMessage,omitempty"`
	IntegerMessage    *uint                 `json:"integerMessage,omitempty"`
}

type receivedInitiatingMessageEventJson struct {
	Model                          string                  `json:"model"`
	CurrentInstance                string                  `json:"currentInstance"`
	Transition                     string                  `json:"transition"`
	InitiatingMessage              string                  `json:"initiatingMessage"`
	NextInstance                   instance.InstanceJson   `json:"nextInstance"`
	RespondingMessage              *message.MessageJson    `json:"respondingMessage,omitempty"`
	RespondingParticipantSignature signature.SignatureJson `json:"respondingParticipantSignature"`
}

type proveMessageExchangeCommandJson struct {
	Model                          string                  `json:"model"`
	CurrentInstance                string                  `json:"currentInstance"`
	Transition                     string                  `json:"transition"`
	InitiatingMessage              string                  `json:"initiatingMessage"`
	Identity                       uint                    `json:"identity"`
	NextInstance                   instance.InstanceJson   `json:"nextInstance"`
	RespondingMessage              *message.MessageJson    `json:"respondingMessage,omitempty"`
	RespondingParticipantSignature signature.SignatureJson `json:"respondingParticipantSignature"`
}

type provedMessageExchangeEventJson struct {
	Instance instance.InstanceJson `json:"instance"`
	Proof    prover.ProofJson      `json:"proof"`
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

func (cmd *proveTerminationCommandJson) ToExecutionCommand() (ProveTerminationCommand, error) {
	return ProveTerminationCommand{
		Model:    cmd.Model,
		Instance: cmd.Instance,
		Identity: cmd.Identity,
	}, nil
}

func (cmd *createInitiatingMessageCommandJson) ToExecutionCommand() (CreateInitiatingMessageCommand, error) {
	var bytesMessage []byte = nil
	var integerMessage *domain.IntegerType = nil
	if cmd.BytesMessage != nil {
		bytesMessage = cmd.BytesMessage
	} else {
		tmp := domain.IntegerType(*cmd.IntegerMessage)
		integerMessage = &tmp
	}
	return CreateInitiatingMessageCommand{
		Model:          cmd.Model,
		Instance:       cmd.Instance,
		Transition:     cmd.Transition,
		BytesMessage:   bytesMessage,
		IntegerMessage: integerMessage,
	}, nil
}

func (cmd *receiveInitiatingMessageCommandJson) ToExecutionCommand() (ReceiveInitiatingMessageCommand, error) {
	model, err := cmd.Model.ToModel()
	if err != nil {
		return ReceiveInitiatingMessageCommand{}, err
	}
	currentInstance, err := cmd.CurrentInstance.ToInstance()
	if err != nil {
		return ReceiveInitiatingMessageCommand{}, err
	}
	initiatingMessage, err := cmd.InitiatingMessage.ToMessage()
	if err != nil {
		return ReceiveInitiatingMessageCommand{}, err
	}
	var bytesMessage []byte = nil
	var integerMessage *domain.IntegerType = nil
	if cmd.BytesMessage != nil {
		bytesMessage = cmd.BytesMessage
	} else {
		tmp := domain.IntegerType(*cmd.IntegerMessage)
		integerMessage = &tmp
	}
	return ReceiveInitiatingMessageCommand{
		Model:             model,
		CurrentInstance:   currentInstance,
		Transition:        cmd.Transition,
		Identity:          cmd.Identity,
		InitiatingMessage: initiatingMessage,
		IntegerMessage:    integerMessage,
		BytesMessage:      bytesMessage,
	}, nil
}

func (cmd *proveMessageExchangeCommandJson) ToExecutionCommand() (ProveMessageExchangeCommand, error) {
	nextInstance, err := cmd.NextInstance.ToInstance()
	if err != nil {
		return ProveMessageExchangeCommand{}, err
	}
	var respondingMessage *domain.Message
	if cmd.RespondingMessage != nil {
		tmp, err := cmd.RespondingMessage.ToMessage()
		if err != nil {
			return ProveMessageExchangeCommand{}, err
		}
		respondingMessage = &tmp
	}
	respondingParticipantSignature, err := cmd.RespondingParticipantSignature.ToSignature()
	if err != nil {
		return ProveMessageExchangeCommand{}, err
	}
	return ProveMessageExchangeCommand{
		Model:                          cmd.Model,
		CurrentInstance:                cmd.CurrentInstance,
		Transition:                     cmd.Transition,
		Identity:                       cmd.Identity,
		InitiatingMessage:              cmd.InitiatingMessage,
		NextInstance:                   nextInstance,
		RespondingMessage:              respondingMessage,
		RespondingParticipantSignature: respondingParticipantSignature,
	}, nil
}

func InstatiatedModelEventToJson(event InstantiatedModelEvent) instantiatedModelEventJson {
	return instantiatedModelEventJson{
		Instance: instance.ToJson(event.Instance),
		Proof:    event.Proof.ToJson(),
	}
}

func ExecutedTransitionEventToJson(event ExecutedTransitionEvent) executedTransitionEventJson {
	return executedTransitionEventJson{
		Instance: instance.ToJson(event.Instance),
		Proof:    event.Proof.ToJson(),
	}
}

func TerminatedInstanceEventToJson(event ProvedTerminationEvent) provedTerminationEventJson {
	return provedTerminationEventJson{
		Proof: event.Proof.ToJson(),
	}
}

func CreatedInitiatingMessageEventToJson(event CreatedInitiatingMessageEvent) createdInitiatingMessageEventJson {
	return createdInitiatingMessageEventJson{
		Model:             model.ToJson(event.Model),
		CurrentInstance:   instance.ToJson(event.CurrentInstance),
		Transition:        event.Transition,
		InitiatingMessage: message.ToJson(event.InintiatingMessage),
	}
}

func ReceivedInitiatingMessageEventToJson(event ReceivedInitiatingMessageEvent) receivedInitiatingMessageEventJson {
	var respondingMessage *message.MessageJson
	if event.RespondingMessage != nil {
		tmp := message.ToJson(*event.RespondingMessage)
		respondingMessage = &tmp
	}
	return receivedInitiatingMessageEventJson{
		Model:                          event.Model,
		CurrentInstance:                event.CurrentInstance,
		Transition:                     event.Transition,
		InitiatingMessage:              event.InitiatingMessage,
		NextInstance:                   instance.ToJson(event.NextInstance),
		RespondingMessage:              respondingMessage,
		RespondingParticipantSignature: signature.ToJson(event.RespondingParticipantSignature),
	}
}

func ProvedMessageExchangeEventToJson(event ProvedMessageExchangeEvent) provedMessageExchangeEventJson {
	return provedMessageExchangeEventJson{
		Instance: instance.ToJson(event.Instance),
		Proof:    event.Proof.ToJson(),
	}
}
