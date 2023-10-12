package state

import "execution-service/domain"

type State struct {
	Model    *domain.Model
	Instance *domain.Instance
	Message  *domain.Message
}
