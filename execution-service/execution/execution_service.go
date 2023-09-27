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
	modelService     domain.ModelService
}

var executionService ExecutionService

func NewExecutionService() ExecutionService {
	if !executionService.isLoaded {
		executionService = ExecutionService{
			isLoaded:        true,
			proofService:    proof.NewProofService(),
			instanceService: domain.NewInstanceService(),
			modelService:    domain.NewModelService(),
		}
	}
	return executionService
}

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (InstantiationResult, error) {
	model := cmd.Model
	instanceResult, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return InstantiationResult{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveInstantiation(model, instanceResult, signature)
	if err != nil {
		return InstantiationResult{}, err
	}
	service.modelService.SaveModel(model)
	service.instanceService.SaveInstance(instanceResult)
	return InstantiationResult{
		Instance: instanceResult,
		Proof:    proofResult,
	}, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (TransitionResult, error) {
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return TransitionResult{}, err
	}
	instance, err := service.instanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return TransitionResult{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return TransitionResult{}, err
	}
	var instanceResult domain.Instance
	if len(cmd.Message) == 0 {
		instanceResult, err = instance.ExecuteTransition(transition)
	} else {
		instanceResult, err = instance.ExecuteTransitionWithMessage(transition, cmd.Message)
	}
	if err != nil {
		return TransitionResult{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveTransition(model, instance, instanceResult, signature)
	if err != nil {
		return TransitionResult{}, err
	}
	service.instanceService.SaveInstance(instanceResult)
	return TransitionResult{
		Instance: instanceResult,
		Proof:    proofResult,
	}, nil
}

func (service *ExecutionService) ProveTermination(cmd ProveTerminationCommand) (TerminationResult, error) {
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return TerminationResult{}, err
	}
	instance, err := service.instanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return TerminationResult{}, err
	}
	signature := service.signatureService.Sign(instance)
	proofResult, err := service.proofService.ProveTermination(model, instance, signature)
	if err != nil {
		return TerminationResult{}, err
	}
	return TerminationResult{
		Proof: proofResult,
	}, nil
}
