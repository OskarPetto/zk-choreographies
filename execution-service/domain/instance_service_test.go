package domain_test

import (
	"proof-service/domain"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindInstanceById(t *testing.T) {
	service := domain.NewInstanceService()
	publicKeys := testdata.GetPublicKeys(1)
	instance := testdata.GetModel2Instance1(publicKeys)
	service.SaveInstance(instance)
	_, err := service.FindInstanceById(instance.Id)
	assert.Nil(t, err)
}

func TestFindInstanceByModel(t *testing.T) {
	service := domain.NewInstanceService()
	publicKeys := testdata.GetPublicKeys(1)
	instance := testdata.GetModel2Instance1(publicKeys)
	service.SaveInstance(instance)
	result := service.FindInstancesByModel(instance.Model)
	assert.Equal(t, 1, len(result))
}
