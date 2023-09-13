package randomness_test

import (
	"proof-service/randomness"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

var instance1 = testdata.GetInstance1()

func TestCreateRandomness(t *testing.T) {
	randomnessService := randomness.NewRandomnessService()
	randomness1, err := randomnessService.CreateRandomness(instance1.Id)
	assert.Nil(t, err)
	randomness2, err := randomnessService.FindRandomness(instance1.Id)
	assert.Equal(t, randomness1, randomness2)
}
