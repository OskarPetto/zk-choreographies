package instance_test

import (
	"execution-service/instance"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstantiateModel(t *testing.T) {
	instanceService := instance.InstanceService{}
	model1 := testdata.GetModel1()
	instance1 := testdata.GetInstance1()
	result, _ := instanceService.InstantiateModel(model1)
	assert.Equal(t, instance1.TokenCounts, result.TokenCounts)
}
