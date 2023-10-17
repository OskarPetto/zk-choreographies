package execution

import (
	"execution-service/domain"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"execution-service/prover"
	"execution-service/signature"
)

type ExecutionService struct {
	ModelService        model.ModelService
	InstanceService     instance.InstanceService
	MessageService      message.MessageService
	ProverService       prover.IProverService
	SignatureService    signature.SignatureService
	SignatureParameters parameters.SignatureParameters
}

func InitializeExecutionService(proverService prover.IProverService) ExecutionService {
	signatureParameters := parameters.NewSignatureParameters()
	modelService := model.NewModelService()
	instanceService := instance.NewInstanceService(modelService)
	messageService := message.NewMessageService(modelService, instanceService, signatureParameters)
	signatureService := signature.NewSignatureService(instanceService)
	return NewExecutionService(modelService, instanceService, messageService, proverService, signatureParameters, signatureService)
}

func NewExecutionService(modelService model.ModelService, instanceService instance.InstanceService, messageService message.MessageService, proverService prover.IProverService, signatureParameters parameters.SignatureParameters, signatureService signature.SignatureService) ExecutionService {
	return ExecutionService{
		ModelService:        modelService,
		InstanceService:     instanceService,
		MessageService:      messageService,
		ProverService:       proverService,
		SignatureParameters: signatureParameters,
		SignatureService:    signatureService,
	}
}

func (service *ExecutionService) InstantiateModel(cmd InstantiateModelCommand) (ExecutionResult, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ExecutionResult{}, err
	}
	instance, err := model.Instantiate(cmd.PublicKeys)
	if err != nil {
		return ExecutionResult{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveInstantiation(prover.ProveInstantiationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return ExecutionResult{}, err
	}
	service.InstanceService.ImportInstance(instance)
	return ExecutionResult{
		Proof:    proof,
		Instance: instance,
	}, nil
}

func (service *ExecutionService) ExecuteTransition(cmd ExecuteTransitionCommand) (ExecutionResult, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ExecutionResult{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ExecutionResult{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return ExecutionResult{}, err
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return ExecutionResult{}, err
	}
	nextInstance, err := currentInstance.ExecuteTransition(transition, constraintInput)
	if err != nil {
		return ExecutionResult{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	senderSignature := nextInstance.Sign(privateKey)
	var recipientSignature *domain.Signature = nil
	if transition.Recipient != domain.EmptyParticipantId {
		tmp, err := service.SignatureService.FindSignatureByInstance(nextInstance.Id()) TODO: split execution and proof
		if err != nil {
			return ExecutionResult{}, err
		}
		recipientSignature = &tmp
	}
	proof, err := service.ProverService.ProveTransition(prover.ProveTransitionCommand{
		Model:              model,
		CurrentInstance:    currentInstance,
		NextInstance:       nextInstance,
		Transition:         transition,
		SenderSignature:    senderSignature,
		RecipientSignature: recipientSignature,
		ConstraintInput:    constraintInput,
	})
	if err != nil {
		return ExecutionResult{}, err
	}
	err = service.InstanceService.ImportInstance(nextInstance)
	if err != nil {
		return ExecutionResult{}, err
	}
	return ExecutionResult{
		Proof:    proof,
		Instance: nextInstance,
	}, nil
}

func (service *ExecutionService) TerminateInstance(cmd TerminateInstanceCommand) (ExecutionResult, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return ExecutionResult{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return ExecutionResult{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	proof, err := service.ProverService.ProveTermination(prover.ProveTerminationCommand{
		Model:     model,
		Instance:  instance,
		Signature: instance.Sign(privateKey),
	})
	if err != nil {
		return ExecutionResult{}, err
	}
	return ExecutionResult{
		Proof:    proof,
		Instance: instance,
	}, nil
}
