package proof_test

import (
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
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance1(publicKeys)
	model := testdata.GetModel2()
	hashService := domain.NewHashService()
	instanceService := domain.NewInstanceService()
	hashService.SaveModelHash(model.Id, domain.HashModel(model))
	instanceService.SaveInstance(instance)
	proofService := proof.NewProofService()

	proof, err := proofService.ProveInstantiation(proof.ProveInstantiationCommand{
		Model:    model,
		Instance: instance.Id(),
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(proof.PublicInput))
}

func TestProveTransition1(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(2)
	currentInstance := testdata.GetModel2Instance2(publicKeys)
	nextInstance := testdata.GetModel2Instance3(publicKeys)
	model := testdata.GetModel2()
	hashService := domain.NewHashService()
	instanceService := domain.NewInstanceService()
	hashService.SaveModelHash(model.Id, domain.HashModel(model))
	instanceService.SaveInstance(currentInstance)
	instanceService.SaveInstance(nextInstance)
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTransition(proof.ProveTransitionCommand{
		Model:           model,
		CurrentInstance: currentInstance.Id(),
		NextInstance:    nextInstance.Id(),
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, len(proof.PublicInput))
}

func TestProveTermination(t *testing.T) {
	publicKeys := testdata.GetPublicKeys(2)
	instance := testdata.GetModel2Instance4(publicKeys)
	model := testdata.GetModel2()
	hashService := domain.NewHashService()
	instanceService := domain.NewInstanceService()
	hashService.SaveModelHash(model.Id, domain.HashModel(model))
	instanceService.SaveInstance(instance)
	proofService := proof.NewProofService()

	proof, err := proofService.ProveTermination(proof.ProveTerminationCommand{
		Model:    model,
		Instance: instance.Id(),
	})
	assert.Nil(t, err)
	assert.Equal(t, 2, len(proof.PublicInput))
}
