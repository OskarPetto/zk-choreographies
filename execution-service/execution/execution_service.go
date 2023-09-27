package execution

import (
	"execution-service/authentication"
	"execution-service/domain"
	"execution-service/proof"
)

type ExecutionService struct {
	isLoaded         bool
	proofService     proof.ProofService
	signatureService authentication.SignatureService
	instanceService  domain.InstanceService
	hashService      domain.HashService
}

var executionService ExecutionService

func NewExecutionService() ExecutionService {
	if !executionService.isLoaded {
		executionService = ExecutionService{
			isLoaded:        true,
			proofService:    proof.NewProofService(),
			instanceService: domain.NewInstanceService(),
			hashService:     domain.NewHashService(),
		}
	}
	return executionService
}

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (InstantiateModelResult, error) {
	model := cmd.Model
	modelHash := domain.HashModel(model)
	instanceResult, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return InstantiateModelResult{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveInstantiation(proof.ProveInstantiationCommand{
		ModelHash: modelHash,
		Model:     model,
		Instance:  instanceResult,
		Signature: signature,
	})
	if err != nil {
		return InstantiateModelResult{}, err
	}
	service.hashService.SaveModelHash(model.Id, modelHash)
	service.instanceService.SaveInstance(instanceResult)
	return InstantiateModelResult{
		Instance: instanceResult,
		Proof:    proofResult,
	}, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (ExecuteTransitionResult, error) {
	model := cmd.Model
	modelHash, err := service.hashService.FindHashByModelId(model.Id)
	if err != nil {
		return ExecuteTransitionResult{}, err
	}
	currentInstance, err := service.instanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ExecuteTransitionResult{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ExecuteTransitionResult{}, err
	}
	var nextInstance domain.Instance
	if len(cmd.Message) == 0 {
		nextInstance, err = currentInstance.ExecuteTransition(transition)
	} else {
		nextInstance, err = currentInstance.ExecuteTransitionWithMessage(transition, cmd.Message)
	}
	if err != nil {
		return ExecuteTransitionResult{}, err
	}
	nextSignature := service.signatureService.Sign(nextInstance)
	proofResult, err := service.proofService.ProveTransition(proof.ProveTransitionCommand{
		ModelHash:       modelHash,
		Model:           model,
		CurrentInstance: currentInstance,
		NextInstance:    nextInstance,
		NextSignature:   nextSignature,
	})
	if err != nil {
		return ExecuteTransitionResult{}, err
	}
	service.instanceService.SaveInstance(nextInstance)
	return ExecuteTransitionResult{
		Instance: nextInstance,
		Proof:    proofResult,
	}, nil
}

func (service *ExecutionService) ProveTermination(cmd ProveTerminationCommand) (ProveTerminationResult, error) {
	model := cmd.Model
	modelHash, err := service.hashService.FindHashByModelId(model.Id)
	if err != nil {
		return ProveTerminationResult{}, err
	}
	instance, err := service.instanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ProveTerminationResult{}, err
	}
	signature := service.signatureService.Sign(instance)
	proofResult, err := service.proofService.ProveTermination(proof.ProveTerminationCommand{
		ModelHash: modelHash,
		Model:     model,
		Instance:  instance,
		Signature: signature,
	})
	if err != nil {
		return ProveTerminationResult{}, err
	}
	return ProveTerminationResult{
		Proof: proofResult,
	}, nil
}
