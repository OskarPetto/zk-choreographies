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

	verifierInput, err := proofService.ProveInstantiation(model, instance, signature)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(verifierInput.Input))
}

func TestProveTransition1(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextInstance.ComputeHash()
	signature := signatureService.Sign(nextInstance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	verifierInput, err := proofService.ProveTransition(model, currentInstance, nextInstance, signature)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(verifierInput.Input))
}

func TestProveTermination(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	verifierInput, err := proofService.ProveTermination(model, instance, signature)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(verifierInput.Input))
}
