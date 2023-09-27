package domain_test

import (
	"execution-service/domain"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindModelById(t *testing.T) {
	service := domain.NewModelService()
	model := testdata.GetModel2()
	service.SaveModel(model)
	_, err := service.FindModelById(model.Id)
	assert.Nil(t, err)
}

func TestFindAllModels(t *testing.T) {
	service := domain.NewModelService()
	model := testdata.GetModel2()
	service.SaveModel(model)
	result := service.FindAllModels()
	assert.Equal(t, 1, len(result))
}
