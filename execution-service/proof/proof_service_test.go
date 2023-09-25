package proof_test

import (
	"proof-service/authentication"
	"proof-service/proof"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProofService(t *testing.T) {
	proof.NewProofService()
}

func TestProveInstantiation(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance1(publicKeys)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	_, err := proofService.ProveInstantiation(instance, model, signature)
	assert.Nil(t, err)
}

func TestProveTransition(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance1(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance2(publicKeys)
	nextInstance.ComputeHash()
	signature := signatureService.Sign(nextInstance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	_, err := proofService.ProveTransition(currentInstance, nextInstance, model, signature)
	assert.Nil(t, err)
}

func TestProveTermination(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	_, err := proofService.ProveTermination(instance, model, signature)
	assert.Nil(t, err)
}
