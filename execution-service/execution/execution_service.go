package execution

import (
	"proof-service/authentication"
	"proof-service/instance"
	"proof-service/model"
	"proof-service/proof"
)

type ExecutionService struct {
	isLoaded         bool
	proofService     proof.ProofService
	signatureService authentication.SignatureService
	instanceService  instance.InstanceService
}

var executionService ExecutionService

func NewExecutionService() ExecutionService {
	if !executionService.isLoaded {
		executionService = ExecutionService{
			isLoaded:        true,
			proofService:    proof.NewProofService(),
			instanceService: instance.NewInstanceService(),
		}
	}
	return executionService
}

func (service *ExecutionService) InstantiatePetriNet(petriNet model.PetriNet, publicKeys [][]byte) (instance.Instance, proof.Proof, error) {
	instanceResult, err := instance.InstantiatePetriNet(petriNet, publicKeys)
	if err != nil {
		return instance.Instance{}, proof.Proof{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveInstantiation(instanceResult, petriNet, signature)
	if err != nil {
		return instance.Instance{}, proof.Proof{}, err
	}
	service.instanceService.SaveInstance(instanceResult)
	return instanceResult, proofResult, nil
}

func (service *ExecutionService) ExecuteTransition(inst instance.Instance, petriNet model.PetriNet, transition model.Transition) (instance.Instance, proof.Proof, error) {
	instanceResult, err := inst.ExecuteTransition(transition)
	if err != nil {
		return instance.Instance{}, proof.Proof{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveTransition(inst, instanceResult, petriNet, signature)
	if err != nil {
		return instance.Instance{}, proof.Proof{}, err
	}
	service.instanceService.SaveInstance(instanceResult)
	return instanceResult, proofResult, nil
}

func (service *ExecutionService) TerminateInstance(inst instance.Instance, petriNet model.PetriNet) (proof.Proof, error) {
	signature := service.signatureService.Sign(inst)
	proofResult, err := service.proofService.ProveTermination(inst, petriNet, signature)
	if err != nil {
		return proof.Proof{}, err
	}
	return proofResult, nil
}
