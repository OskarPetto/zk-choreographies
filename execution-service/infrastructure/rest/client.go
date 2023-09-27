package rest

import (
	"execution-service/domain"
	"execution-service/infrastructure/json"
	"io"
	"net/http"
	"os"
)

type ModelServiceClient struct {
	modelServiceUrl string
}

func NewModelClient() domain.ModelService {
	return ModelServiceClient{
		modelServiceUrl: os.Getenv("MODEL_SERVICE_URL"),
	}
}

func (client ModelServiceClient) FindModelById(modelId domain.ModelId) (domain.Model, error) {

	response, err := http.Get(client.modelServiceUrl + "/models/" + modelId)
	if err != nil {
		return domain.Model{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.Model{}, err
	}
	model, err := json.UnmarshalModel(body)
	if err != nil {
		return domain.Model{}, err
	}
	return model, nil
}
