package instance_test

import (
	"execution-service/instance"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindModel(t *testing.T) {
	instanceService := instance.NewInstanceService()
	instance1 := testdata.GetInstance1()
	instanceService.SaveInstance(instance1)
	result, _ := instanceService.FindInstance(instance1.Id)
	assert.Equal(t, instance1, result)
}
