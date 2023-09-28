package model_test

import (
	"execution-service/domain"
	"execution-service/model"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ModelPortMock struct {
}

func (service ModelPortMock) FindModelById(id domain.ModelId) (domain.Model, error) {
	return testdata.GetModel2(), nil
}

var modelService = model.NewModelService(ModelPortMock{})

func TestFindModelById(t *testing.T) {
	model := testdata.GetModel2()
	modelResult, err := modelService.FindModelById(model.Id)
	assert.Nil(t, err)
	assert.NotEqual(t, model.Hash, modelResult.Hash)
}
