package authentication_test

import (
	"proof-service/authentication"
	"proof-service/testdata"
	"testing"
)

func TestSign(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	instance := testdata.GetPetriNet1Instance1(testdata.GetPublicKeys(1)[0])
	signatureService.Sign(instance)
}
