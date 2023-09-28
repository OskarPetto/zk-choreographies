package instance_test

import (
	"execution-service/instance"
	"execution-service/signature"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureService signature.SignatureService = signature.InitializeSignatureService()

func TestFindInstancesByModel(t *testing.T) {
	service := instance.NewInstanceService()
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	service.ImportInstance(instance)
	result := service.FindInstancesByModel(instance.Model)
	assert.Equal(t, 1, len(result))
}

func TestFindInstanceById(t *testing.T) {
	service := instance.NewInstanceService()
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	service.ImportInstance(instance)
	_, err := service.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}
