package state_test

import (
	"execution-service/domain"
	"execution-service/state"
	"testing"

	"github.com/stretchr/testify/assert"
)

var service = state.InitializeStateService()

func TestImportState(t *testing.T) {
	privateKey1 := signatureParameters.GetPrivateKeyForIdentity(0)
	privateKey2 := signatureParameters.GetPrivateKeyForIdentity(1)
	publicKey2 := domain.NewPublicKey(privateKey2.PublicKey)
	model := &states[0].Model
	instance := states[0].Instance
	plaintext := state.SerializeInstance(instance)
	encryptedInstance := plaintext.Encrypt(privateKey1, publicKey2)
	domainState := domain.State{
		Model:             model,
		EncryptedInstance: &encryptedInstance,
		Message:           nil,
	}
	cmd := state.ImportStateCommand{
		State:    domainState,
		Identity: 1,
	}
	err := service.ImportState(cmd)
	assert.Nil(t, err)
	_, err = service.ModelService.FindModelById(model.Id())
	assert.Nil(t, err)
	_, err = service.InstanceService.FindInstanceById(instance.Id())
	assert.Nil(t, err)
}
