package domain_test

import (
	"proof-service/domain"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindInstance(t *testing.T) {
	service := domain.NewInstanceService()
	publicKeys := testdata.GetPublicKeys(1)
	instance := testdata.GetModel2Instance1(publicKeys)
	service.SaveInstance(instance)
	_, err := service.FindInstanceByHash(instance.Hash)
	assert.Nil(t, err)
}
