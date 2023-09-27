package domain

type ModelService interface {
	FindModelById(ModelId) Model
}
