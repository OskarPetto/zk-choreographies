package testdata

import (
	"execution-service/domain"
	"execution-service/infrastructure/json"
	"execution-service/utils"
	"io"
	"os"
)

func GetModel2() domain.Model {
	jsonFile, err := os.Open("/home/opetto/uni/zk-choreographies/execution-service/testdata/example_choreography.json")
	utils.PanicOnError(err)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	model, err := json.UnmarshalModel(byteValue)
	utils.PanicOnError(err)
	return model
}
