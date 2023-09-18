package workflow_test

import (
	"proof-service/testdata"
	"proof-service/workflow"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeInstance(t *testing.T) {
	inst := testdata.GetPetriNet1Instance1()
	expected := testdata.GetPetriNet1Instance1Serialized()
	result := workflow.SerializeInstance(inst)
	assert.Equal(t, len(expected), len(result))
	//assert.Equal(t, expected, result)
}
