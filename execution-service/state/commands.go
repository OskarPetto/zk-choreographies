package state

import "execution-service/domain"

type ImportStateCommand struct {
	EncryptedState domain.Ciphertext
	Identity       domain.IdentityId
}
