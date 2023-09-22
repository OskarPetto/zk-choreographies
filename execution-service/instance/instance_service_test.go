package instance_test

import (
	"proof-service/instance"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindInstance(t *testing.T) {
	service := instance.NewInstanceService()
	publicKey := testdata.GetPublicKeys(1)[0]
	instance := testdata.GetPetriNet1Instance1(publicKey)
	service.SaveInstance(instance)
	_, err := service.FindInstanceByHash(instance.Hash)
	assert.Nil(t, err)
}
