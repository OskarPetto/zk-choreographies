package testdata

import (
	"io"
	"os"
	"proof-service/domain"
	"proof-service/infrastructure"
	"proof-service/utils"
)

func GetModel2() domain.Model {
	jsonFile, err := os.Open("/home/opetto/uni/zk-choreographies/execution-service/testdata/example_choreography.json")
	utils.PanicOnError(err)
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	return infrastructure.FromJson(byteValue)
}
