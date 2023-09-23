package execution

import (
	"proof-service/authentication"
	"proof-service/domain"
	"proof-service/proof"
)

type ExecutionService struct {
	isLoaded         bool
	proofService     proof.ProofService
	signatureService authentication.SignatureService
	instanceService  domain.InstanceService
}

var executionService ExecutionService

func NewExecutionService() ExecutionService {
	if !executionService.isLoaded {
		executionService = ExecutionService{
			isLoaded:        true,
			proofService:    proof.NewProofService(),
			instanceService: domain.NewInstanceService(),
		}
	}
	return executionService
}

func (service *ExecutionService) InstantiateModel(model domain.Model, publicKeys [][]byte) (domain.Instance, proof.Proof, error) {
	instanceResult, err := model.Instantiate(publicKeys)
	if err != nil {
		return domain.Instance{}, proof.Proof{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveInstantiation(instanceResult, model, signature)
	if err != nil {
		return domain.Instance{}, proof.Proof{}, err
	}
	service.instanceService.SaveInstance(instanceResult)
	return instanceResult, proofResult, nil
}

func (service *ExecutionService) ExecuteTransition(inst domain.Instance, model domain.Model, transition domain.Transition) (domain.Instance, proof.Proof, error) {
	instanceResult, err := inst.ExecuteTransition(transition)
	if err != nil {
		return domain.Instance{}, proof.Proof{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveTransition(inst, instanceResult, model, signature)
	if err != nil {
		return domain.Instance{}, proof.Proof{}, err
	}
	service.instanceService.SaveInstance(instanceResult)
	return instanceResult, proofResult, nil
}

func (service *ExecutionService) TerminateInstance(inst domain.Instance, model domain.Model) (proof.Proof, error) {
	signature := service.signatureService.Sign(inst)
	proofResult, err := service.proofService.ProveTermination(inst, model, signature)
	if err != nil {
		return proof.Proof{}, err
	}
	return proofResult, nil
}
