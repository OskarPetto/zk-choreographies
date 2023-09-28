package signature_test

import (
	"execution-service/signature"
	"execution-service/testdata"
	"testing"
)

var signatureService signature.SignatureService = signature.InitializeSignatureService()

func TestSign(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	signatureService.Sign(instance)
}
