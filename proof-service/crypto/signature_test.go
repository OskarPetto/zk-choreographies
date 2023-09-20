package crypto_test

import (
	"proof-service/crypto"
	"proof-service/testdata"
	"testing"
)

func TestSign(t *testing.T) {
	crypto.LoadParameters()
	commitment := testdata.GetCommitment1()
	crypto.Sign(commitment)
}
