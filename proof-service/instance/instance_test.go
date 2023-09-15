package instance_test

import (
	"proof-service/instance"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInstance(t *testing.T) {
	inst := testdata.GetInstance1()
	err := instance.ValidateInstance(inst)
	assert.Nil(t, err)
}

func TestSerializeInstance(t *testing.T) {
	inst := testdata.GetInstance1()
	expected := testdata.GetSerializedInstance1()
	result := instance.SerializeInstance(inst)
	assert.Equal(t, expected, result)
}
