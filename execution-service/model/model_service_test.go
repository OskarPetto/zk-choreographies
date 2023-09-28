package model_test

import (
	"execution-service/model"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var modelService = model.NewModelService()

func TestImportModel(t *testing.T) {
	model := testdata.GetModel2()
	modelService.ImportModel(model)
	modelResult, err := modelService.FindModelById(model.Id())
	assert.Nil(t, err)
	assert.Equal(t, model.Hash, modelResult.Hash)
}

func TestCreateModel(t *testing.T) {
	model := testdata.GetModel2()
	modelResult := modelService.CreateModel(model)
	assert.NotEqual(t, model.Hash, modelResult.Hash)
}
