package domain_test

import (
	"execution-service/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecryptState(t *testing.T) {
	privateKey1 := signatureParameters.GetPrivateKeyForIdentity(0)
	privateKey2 := signatureParameters.GetPrivateKeyForIdentity(1)
	publicKey2 := domain.NewPublicKey(privateKey2.PublicKey)
	state := domain.Plaintext{Value: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}
	ciphertext := state.Encrypt(privateKey1, publicKey2)
	result, err := ciphertext.Decrypt(privateKey2)
	assert.Nil(t, err)
	assert.Equal(t, state, result)
}
