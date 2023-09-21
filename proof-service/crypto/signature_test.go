package crypto_test

import (
	"proof-service/crypto"
	"proof-service/testdata"
	"testing"
)

func TestSign(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	commitment := testdata.GetCommitment1()
	signatureService.Sign(commitment)
}
