package state

import "execution-service/domain"

type ImportStateCommand struct {
	EncryptedState domain.EncryptedState
	Identity       domain.IdentityId
}
