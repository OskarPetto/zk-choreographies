package instance_test

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/testdata"
	"execution-service/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters parameters.SignatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)

var modelService = model.NewModelService()
var instanceService = instance.NewInstanceService(modelService)

func TestFindInstancesByModel(t *testing.T) {
	instance := states[0].Instance
	instanceService.SaveInstance(instance)
	modelId := utils.BytesToString(instance.Model.Value[:])
	result := instanceService.FindInstancesByModel(modelId)
	assert.Equal(t, 1, len(result))
}

func TestFindInstanceById(t *testing.T) {
	instance := states[0].Instance
	instanceService.SaveInstance(instance)
	_, err := instanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}

func TestImportInstanceWithoutModel(t *testing.T) {
	instance := states[0].Instance
	err := instanceService.ImportInstance(instance)
	assert.NotNil(t, err)
}

func TestImportInstanceWithModel(t *testing.T) {
	model := states[0].Model
	modelService.ImportModel(model)
	instance := states[0].Instance
	err := instanceService.ImportInstance(instance)
	assert.Nil(t, err)
}

func TestImportInstanceInvalidHash(t *testing.T) {
	instance := states[0].Instance
	instance.Hash = domain.SaltedHash{}
	err := instanceService.ImportInstance(instance)
	assert.NotNil(t, err)
}
