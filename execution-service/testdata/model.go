package testdata

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/model"
	"execution-service/utils"
	"io"
	"os"
)

func fromJson(data []byte) (domain.Model, error) {
	var model model.ModelJson
	err := json.Unmarshal(data, &model)
	if err != nil {
		return domain.Model{}, err
	}
	return model.ToModel()
}

func GetModel2() domain.Model {
	jsonFile, err := os.Open("../testdata/example_choreography.json")
	utils.PanicOnError(err)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	model, err := fromJson(byteValue)
	utils.PanicOnError(err)
	model.CreatedAt = 1695983320
	model.UpdateHash()
	return model
}
