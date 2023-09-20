package crypto_test

import (
	"fmt"
	"proof-service/crypto"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommitment(t *testing.T) {
	instance := testdata.GetPetriNet1Instance1()
	result := crypto.NewCommitment(instance)
	assert.Equal(t, instance.Id, result.Id)
	fmt.Printf("%+v\n", result)
}
