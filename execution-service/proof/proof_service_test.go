package proof_test

import (
	"proof-service/crypto"
	"proof-service/proof"
	"proof-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProofService(t *testing.T) {
	proof.NewProofService()
}

func TestProveInstantiation(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance1(publicKey)
	petriNet := testdata.GetPetriNet1()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveInstantiation(instance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTransition(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	currentInstance := testdata.GetPetriNet1Instance1(publicKey)
	nextInstance := testdata.GetPetriNet1Instance2(publicKey)
	petriNet := testdata.GetPetriNet1()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTransition(currentInstance, nextInstance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}

func TestProveTermination(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance3(publicKey)
	petriNet := testdata.GetPetriNet1()
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTermination(instance, petriNet)
	assert.Nil(t, err)
	assert.Equal(t, 128, len(proof))
}
