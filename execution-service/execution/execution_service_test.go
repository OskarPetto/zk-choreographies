package execution_test

import (
	"execution-service/domain"
	"execution-service/execution"
	"execution-service/testdata"
	"testing"
)

type ModelServiceMock struct {
}

func (service ModelServiceMock) FindModelById(id domain.ModelId) (domain.Model, error) {
	return testdata.GetModel2(), nil
}

func initializeExecutionService() execution.ExecutionService {
	return execution.InitializeExecutionService(ModelServiceMock{})
}

func TestInitializeExecutionService(t *testing.T) {
	initializeExecutionService()
}
