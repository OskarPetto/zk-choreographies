package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
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

type terminateInstanceCommandJson struct {
	Model    string `json:"model"`
	Instance string `json:"instance"`
	Identity uint   `json:"identity"`
}

type terminatedInstanceEventJson struct {
	Proof prover.ProofJson `json:"proof"`
}

type sendMessageCommandJson struct {
	Model          string `json:"model"`
	Instance       string `json:"instance"`
	Transition     string `json:"transition"`
	Identity       uint   `json:"identity"`
	BytesMessage   []byte `json:"bytesMessage,omitempty"`
	IntegerMessage *uint  `json:"integerMessage,omitempty"`
}

type sentMessageEventJson struct {
	Model            string                  `json:"model"`
	CurrentInstance  string                  `json:"currentInstance"`
	Transition       string                  `json:"transition"`
	NextInstance     instance.InstanceJson   `json:"nextInstance"`
	SenderSignature  signature.SignatureJson `json:"senderSignature"`
	EncryptedMessage *message.CiphertextJson `json:"encryptedMessage,omitempty"`
}

type receiveMessageCommandJson struct {
	Model            string                  `json:"model"`
	CurrentInstance  string                  `json:"currentInstance"`
	Transition       string                  `json:"transition"`
	NextInstance     instance.InstanceJson   `json:"nextInstance"`
	SenderSignature  signature.SignatureJson `json:"senderSignature"`
	EncryptedMessage *message.CiphertextJson `json:"encryptedMessage,omitempty"`
	Identity         uint                    `json:"identity"`
}

type receivedMessageEventJson struct {
	Proof prover.ProofJson `json:"proof"`
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

func (cmd *sendMessageCommandJson) ToExecutionCommand() (SendMessageCommand, error) {
	var bytesMessage []byte = nil
	var integerMessage *domain.IntegerType = nil
	if cmd.BytesMessage != nil {
		bytesMessage = cmd.BytesMessage
	} else {
		tmp := domain.IntegerType(*cmd.IntegerMessage)
		integerMessage = &tmp
	}
	return SendMessageCommand{
		Model:          cmd.Model,
		Instance:       cmd.Instance,
		Transition:     cmd.Transition,
		Identity:       cmd.Identity,
		BytesMessage:   bytesMessage,
		IntegerMessage: integerMessage,
	}, nil
}

func (cmd *receiveMessageCommandJson) ToExecutionCommand() (ReceiveMessageCommand, error) {
	nextInstance, err := cmd.NextInstance.ToInstance()
	if err != nil {
		return ReceiveMessageCommand{}, err
	}
	senderSignature, err := cmd.SenderSignature.ToSignature()
	if err != nil {
		return ReceiveMessageCommand{}, err
	}
	var encryptedMessage *domain.Ciphertext
	if cmd.EncryptedMessage != nil {
		tmp, err := cmd.EncryptedMessage.ToCiphertext()
		if err != nil {
			return ReceiveMessageCommand{}, err
		}
		encryptedMessage = &tmp
	}

	return ReceiveMessageCommand{
		Model:            cmd.Model,
		CurrentInstance:  cmd.CurrentInstance,
		Transition:       cmd.Transition,
		Identity:         cmd.Identity,
		NextInstance:     nextInstance,
		SenderSignature:  senderSignature,
		EncryptedMessage: encryptedMessage,
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

func TerminatedInstanceEventToJson(event TerminatedInstanceEvent) terminatedInstanceEventJson {
	return terminatedInstanceEventJson{
		Proof: event.Proof.ToJson(),
	}
}

func SentMessageEventToJson(event SentMessageEvent) sentMessageEventJson {
	var encryptedMessage *message.CiphertextJson
	if event.EncryptedMessage != nil {
		tmp := message.ToCiphertextJson(*event.EncryptedMessage)
		encryptedMessage = &tmp
	}

	return sentMessageEventJson{
		Model:            event.Model,
		CurrentInstance:  event.CurrentInstance,
		Transition:       event.Transition,
		NextInstance:     instance.ToJson(event.NextInstance),
		SenderSignature:  signature.ToJson(event.SenderSignature),
		EncryptedMessage: encryptedMessage,
	}
}

func ReceivedMessageEventToJson(event ReceivedMessageEvent) receivedMessageEventJson {
	return receivedMessageEventJson{
		Proof: event.Proof.ToJson(),
	}
}
