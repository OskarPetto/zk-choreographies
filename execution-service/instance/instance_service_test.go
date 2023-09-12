package instance_test

import (
	"execution-service/instance"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var instance1 = testdata.GetInstance1()
var instanceService = instance.NewInstanceService()

func TestFindInstance(t *testing.T) {
	instanceService.SaveInstance(instance1)
	result, _ := instanceService.FindInstance(instance1.Id)
	assert.Equal(t, instance1, result)
}
