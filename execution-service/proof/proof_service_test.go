package proof_test

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/model"
	"execution-service/proof"
	"execution-service/signature"
	"execution-service/testdata"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ModelServiceMock struct {
}

func (service ModelServiceMock) FindModelById(id domain.ModelId) (domain.Model, error) {
	return testdata.GetModel2(), nil
}

var signatureService signature.SignatureService
var proofService proof.ProofService
var instanceService instance.InstanceService
var modelService model.ModelService

func TestInitializeProofService(t *testing.T) {
	proofService = proof.InitializeProofService()
	instanceService = proofService.InstanceService
	modelService = proofService.ModelService
	signatureService = proofService.SignatureService
}

func TestProveInstantiation(t *testing.T) {
	model := testdata.GetModel2()
	modelService.ImportModel(model)
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance1(publicKeys)
	instanceService.ImportInstance(instance)

	proof, err := proofService.ProveInstantiation(proof.ProveInstantiationCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(proof.PublicInput))
}

func TestProveTransition1(t *testing.T) {
	model := testdata.GetModel2()
	modelService.ImportModel(model)
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	instanceService.ImportInstance(currentInstance)
	instanceService.ImportInstance(nextInstance)

	proof, err := proofService.ProveTransition(proof.ProveTransitionCommand{
		Model:           model.Id(),
		CurrentInstance: currentInstance.Id(),
		NextInstance:    nextInstance.Id(),
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(proof.PublicInput))
}

func TestProveTermination(t *testing.T) {
	model := testdata.GetModel2()
	modelService.ImportModel(model)
	publicKeys := testdata.GetPublicKeys(signatureService, 2)
	instance := testdata.GetModel2Instance4(publicKeys)
	instanceService.ImportInstance(instance)

	proof, err := proofService.ProveTermination(proof.ProveTerminationCommand{
		Model:    model.Id(),
		Instance: instance.Id(),
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(proof.PublicInput))
}
