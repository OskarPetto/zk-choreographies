package execution

import (
	"execution-service/authentication"
	"execution-service/domain"
	"execution-service/proof"
	"fmt"
)

type InstantiateModelCommand struct {
	Model      domain.Model
	PublicKeys []domain.PublicKey
}

type ExecuteTransitionCommand struct {
	Model      domain.ModelId
	Instance   domain.InstanceId
	Transition domain.TransitionId
	Message    []byte
}

type TerminateInstanceCommand struct {
	Model    domain.ModelId
	Instance domain.InstanceId
}

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

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (domain.Instance, error) {
	model := cmd.Model
	instanceResult, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return domain.Instance{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveInstantiation(model, instanceResult, signature)
	if err != nil {
		return domain.Instance{}, err
	}
	//TODO call ethereumservice with proofResult
	fmt.Println(proofResult.PublicInput)
	service.modelService.SaveModel(model)
	service.instanceService.SaveInstance(instanceResult)
	return instanceResult, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (domain.Instance, error) {
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return domain.Instance{}, err
	}
	instance, err := service.instanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return domain.Instance{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return domain.Instance{}, err
	}
	var instanceResult domain.Instance
	if len(cmd.Message) == 0 {
		instanceResult, err = instance.ExecuteTransition(transition)
	} else {
		instanceResult, err = instance.ExecuteTransitionWithMessage(transition, cmd.Message)
	}
	if err != nil {
		return domain.Instance{}, err
	}
	signature := service.signatureService.Sign(instanceResult)
	proofResult, err := service.proofService.ProveTransition(model, instance, instanceResult, signature)
	if err != nil {
		return domain.Instance{}, err
	}
	//TODO call ethereumservice with proofResult
	fmt.Println(proofResult.PublicInput)
	service.instanceService.SaveInstance(instanceResult)
	return instanceResult, nil
}

func (service *ExecutionService) TerminateInstance(cmd TerminateInstanceCommand) error {
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return err
	}
	instance, err := service.instanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return err
	}
	signature := service.signatureService.Sign(instance)
	proofResult, err := service.proofService.ProveTermination(model, instance, signature)
	if err != nil {
		return err
	}
	//TODO call ethereumservice with proofResult
	fmt.Println(proofResult.PublicInput)
	return nil
}
