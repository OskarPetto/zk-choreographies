package authentication_test

import (
	"proof-service/authentication"
	"proof-service/testdata"
	"testing"
)

func TestSign(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	signatureService.Sign(instance)
}
