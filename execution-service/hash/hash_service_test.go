package hash_test

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindHashByModelId(t *testing.T) {
	service := hash.NewHashService()
	model := testdata.GetModel2()
	service.SaveModelHash(model.Id, domain.HashModel(model))
	_, err := service.FindHashByModelId(model.Id)
	assert.Nil(t, err)
}
