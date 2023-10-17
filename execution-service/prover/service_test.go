package prover_test

import (
	"encoding/json"
	"execution-service/parameters"
	"execution-service/prover"
	"execution-service/testdata"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var proofService prover.ProverService
var signatureParameters = parameters.NewSignatureParameters()
var states = testdata.GetModel2States(signatureParameters)
var proofs []prover.ProofJson

func TestInitializeProofService(t *testing.T) {
	proofService = prover.InitializeProverService()
}

func TestProveInstantiation(t *testing.T) {
	instance := states[0].Instance
	model := states[0].Model
	signature := states[0].SenderSignature

	proof, err := proofService.ProveInstantiation(prover.ProveInstantiationCommand{
		Model:     model,
		Instance:  instance,
		Signature: signature,
	})
	assert.Nil(t, err)
	proofs = append(proofs, proof.ToJson())
}

func TestProveTransition0(t *testing.T) {
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	senderSignature := states[1].SenderSignature
	recipientSignature := states[1].RecipientSignature
	constraintInput := states[1].ConstraintInput

	proof, err := proofService.ProveTransition(prover.ProveTransitionCommand{
		Model:              model,
		CurrentInstance:    currentInstance,
		NextInstance:       nextInstance,
		Transition:         model.Transitions[0],
		SenderSignature:    senderSignature,
		RecipientSignature: recipientSignature,
		ConstraintInput:    constraintInput,
	})
	assert.Nil(t, err)
	proofs = append(proofs, proof.ToJson())
}

func TestProveTermination(t *testing.T) {
	instance := states[len(states)-1].Instance
	model := states[len(states)-1].Model
	signature := states[len(states)-1].SenderSignature

	proof, err := proofService.ProveTermination(prover.ProveTerminationCommand{
		Model:     model,
		Instance:  instance,
		Signature: signature,
	})
	assert.Nil(t, err)
	proofs = append(proofs, proof.ToJson())

	bytes, _ := json.Marshal(proofs)
	fmt.Println(string(bytes))
}
