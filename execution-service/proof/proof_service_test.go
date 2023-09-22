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
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance1(publicKey)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	petriNet := testdata.GetPetriNet1()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveInstantiation(instance, petriNet, signature)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTransition(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	currentInstance.ComputeHash()
	nextInstance := testdata.GetPetriNet1Instance2(publicKey)
	nextInstance.ComputeHash()
	signature := signatureService.Sign(nextInstance)
	petriNet := testdata.GetPetriNet1()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTransition(currentInstance, nextInstance, petriNet, signature)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTermination(t *testing.T) {
	signatureService := authentication.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance3(publicKey)
	instance.ComputeHash()
	signature := signatureService.Sign(instance)
	petriNet := testdata.GetPetriNet1()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTermination(instance, petriNet, signature)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}
