package petri_net_test

import (
	"proof-service/petri_net"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePetriNet(t *testing.T) {
	petriNet := testdata.GetPetriNet()
	err := petri_net.ValidatePetriNet(petriNet)
	assert.Nil(t, err)
}
