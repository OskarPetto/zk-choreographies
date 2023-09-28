package model

import (
	"execution-service/domain"
)

type ModelPort interface {
	FindModelById(domain.ModelId) (domain.Model, error)
}
