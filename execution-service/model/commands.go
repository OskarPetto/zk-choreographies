package model

import "execution-service/domain"

type ImportModelCommand struct {
	Model    domain.Model
	Instance domain.Instance
}
