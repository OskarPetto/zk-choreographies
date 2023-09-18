package instance_test

import (
	"proof-service/instance"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromWorkflowInstance(t *testing.T) {
	workflowInstance := testdata.GetWorkflowInstance1()
	expected := testdata.GetPetriNet1Instance1()
	result, err := instance.FromWorkflowInstance(workflowInstance)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestSerializeInstance(t *testing.T) {
	inst := testdata.GetPetriNet1Instance1()
	expected := testdata.GetPetriNet1Instance1Serialized()
	result := instance.SerializeInstance(inst)
	assert.Equal(t, len(expected), len(result))
	//assert.Equal(t, expected, result)
}
