package rest

import (
	"execution-service/domain"
	"net/http"
)

type ModelClient struct {
}

func (client *ModelClient) FindModelById(modelId domain.ModelId) domain.Model {

	response, err := http.Get("http://pokeapi.co/api/v2/pokedex/kanto/")
}
