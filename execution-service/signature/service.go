package signature

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/utils"
	"fmt"
)

type SignatureService struct {
	InstanceService instance.InstanceService
	signatures      map[string]domain.Signature
}

func NewSignatureService(instanceService instance.InstanceService) SignatureService {
	return SignatureService{
		signatures:      make(map[string]domain.Signature),
		InstanceService: instanceService,
	}
}

func (service *SignatureService) FindSignatureByInstance(instanceId domain.InstanceId) (domain.Signature, error) {
	signature, exists := service.signatures[instanceId]
	if !exists {
		return domain.Signature{}, fmt.Errorf("signature for instance %s not found", instanceId)
	}
	return signature, nil
}

func (service *SignatureService) ImportSignature(signature domain.Signature) error {
	instanceId := utils.BytesToString(signature.Instance.Value[:])
	_, err := service.InstanceService.FindInstanceById(instanceId)
	if err != nil {
		return err
	}
	if !signature.Verify() {
		return fmt.Errorf("signature does not verify")
	}
	service.signatures[instanceId] = signature
	return nil
}
