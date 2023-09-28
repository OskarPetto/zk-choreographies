package adapter

import (
	"execution-service/domain"
	"execution-service/model"
	"io"
	"net/http"
	"os"
)

type ModelAdapter struct {
	modelServiceUrl string
}

func NewModelAdapter() *ModelAdapter {
	return &ModelAdapter{
		modelServiceUrl: os.Getenv("MODEL_SERVICE_URL"),
	}
}

func (client ModelAdapter) FindModelById(modelId domain.ModelId) (domain.Model, error) {

	response, err := http.Get(client.modelServiceUrl + "/models/" + modelId)
	if err != nil {
		return domain.Model{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.Model{}, err
	}
	model, err := model.FromJson(body)
	if err != nil {
		return domain.Model{}, err
	}
	return model, nil
}
