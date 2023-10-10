package prover_test

import (
	"encoding/json"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/prover"
	"execution-service/testdata"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var proofService prover.ProverService
var instanceService instance.InstanceService
var modelService model.ModelService
var signatureParameters parameters.SignatureParameters
var states []testdata.State
var proofs []prover.ProofJson

func TestInitializeProofService(t *testing.T) {
	proofService = prover.InitializeProverService()
	instanceService = proofService.InstanceService
	modelService = proofService.ModelService
	signatureParameters = proofService.SignatureParameters
	states = testdata.GetModel2States(signatureParameters)
	for _, state := range states {
		modelService.ImportModel(state.Model)
		instanceService.ImportInstance(state.Instance)
	}
}

func TestProveInstantiation(t *testing.T) {
	instance := states[0].Instance
	model := states[0].Model
	identity := states[0].Identity

	proof, err := proofService.ProveInstantiation(prover.ProveInstantiationCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
		Identity: identity,
	})
	assert.Nil(t, err)
	proofs = append(proofs, proof.ToJson())
}

func TestProveTransition0(t *testing.T) {
	model := states[0].Model
	currentInstance := states[0].Instance
	nextInstance := states[1].Instance
	identity := states[1].Identity

	proof, err := proofService.ProveTransition(prover.ProveTransitionCommand{
		Model:           model.Id(),
		CurrentInstance: currentInstance.Id(),
		NextInstance:    nextInstance.Id(),
		Transition:      model.Transitions[0].Id,
		Identity:        identity,
	})
	assert.Nil(t, err)
	proofs = append(proofs, proof.ToJson())
}

func TestProveTermination(t *testing.T) {
	instance := states[len(states)-1].Instance
	model := states[len(states)-1].Model
	identity := states[len(states)-1].Identity

	proof, err := proofService.ProveTermination(prover.ProveTerminationCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
		Identity: identity,
		EndPlace: 13,
	})
	assert.Nil(t, err)
	proofs = append(proofs, proof.ToJson())

	bytes, _ := json.Marshal(proofs)
	fmt.Println(string(bytes))
}
