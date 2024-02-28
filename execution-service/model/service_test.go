package model_test

import (
	"execution-service/instance"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var signatureParameters = parameters.NewSignatureParameters()
var instanceService = instance.NewInstanceService()
var modelService = model.NewModelService()
var states = testdata.GetModel2States(signatureParameters)

func TestImportModel(t *testing.T) {
	domainModel := states[0].Model
	modelService.ImportModel(domainModel)
	modelResult, err := modelService.FindModelById(domainModel.Id())
	assert.Nil(t, err)
	assert.Equal(t, domainModel.SaltedHash, modelResult.SaltedHash)
}

func TestCreateModel(t *testing.T) {
	model := testdata.GetModel2()
	modelResult := modelService.CreateModel(model)
	assert.NotEqual(t, model.SaltedHash, modelResult)
}
