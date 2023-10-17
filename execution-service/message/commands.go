package message

import "execution-service/domain"

type CreateMessageCommand struct {
	Model          domain.ModelId
	Instance       domain.InstanceId
	Transition     domain.TransitionId
	Identity       domain.IdentityId
	BytesMessage   []byte
	IntegerMessage *domain.IntegerType
}

type CreateMessageResult struct {
	Instance         domain.Instance
	EncryptedMessage domain.Ciphertext
}

type ImportMessageCommand struct {
	Instance         domain.InstanceId
	EncryptedMessage domain.Ciphertext
	Identity         domain.IdentityId
}
