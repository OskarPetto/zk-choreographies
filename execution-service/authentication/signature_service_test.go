package authentication_test

import (
	"execution-service/authentication"
	"execution-service/testdata"
	"testing"
)

func TestSign(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetModel2Instance1(testdata.GetPublicKeys(2))
	signatureService.Sign(instance)
}
