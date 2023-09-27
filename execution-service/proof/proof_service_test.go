package proof_test

import (
	"execution-service/authentication"
	"execution-service/domain"
	"execution-service/proof"
	"execution-service/testdata"
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
	signature := signatureService.Sign(instance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveInstantiation(proof.ProveInstantiationCommand{
		ModelHash: domain.HashModel(model),
		Model:     model,
		Instance:  instance,
		Signature: signature,
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(proof.PublicInput))
}

func TestProveTransition1(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	nextSignature := signatureService.Sign(nextInstance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTransition(proof.ProveTransitionCommand{
		ModelHash:       domain.HashModel(model),
		Model:           model,
		CurrentInstance: currentInstance,
		NextInstance:    nextInstance,
		NextSignature:   nextSignature,
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(proof.PublicInput))
}

func TestProveTermination(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	signature := signatureService.Sign(instance)
	model := testdata.GetModel2()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTermination(proof.ProveTerminationCommand{
		ModelHash: domain.HashModel(model),
		Model:     model,
		Instance:  instance,
		Signature: signature,
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(proof.PublicInput))
}
