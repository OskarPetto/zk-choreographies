package petri_net_test

import (
	"proof-service/petri_net"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromWorkflowPetriNet(t *testing.T) {
	workflowPetriNet := testdata.GetWorkflowPetriNet1()
	expected := testdata.GetPetriNet1()
	result, err := petri_net.FromWorkflowPetriNet(workflowPetriNet)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
