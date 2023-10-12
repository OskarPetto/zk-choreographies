package state

import "execution-service/domain"

type ImportStateCommand struct {
	EncryptedState domain.Chiphertext
	Identity       domain.IdentityId
}
