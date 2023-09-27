package domain

type ModelService interface {
	FindModelById(ModelId) (Model, error)
}

var ModelServiceImpl ModelService
