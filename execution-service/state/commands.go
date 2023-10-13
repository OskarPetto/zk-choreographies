package state

import "execution-service/domain"

type ImportStateCommand struct {
	State    domain.State
	Identity domain.IdentityId
}
