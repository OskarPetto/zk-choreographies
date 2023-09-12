package execution_test

import (
	"execution-service/execution"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstantiateModel(t *testing.T) {
	executionService := execution.ExecutionService{}
	model1 := testdata.GetModel1()
	instance1 := testdata.GetInstance1()
	result := executionService.InstantiateModel(model1)
	assert.Equal(t, instance1.Model, result.Model)
	assert.Equal(t, instance1.TokenCounts, result.TokenCounts)
}
