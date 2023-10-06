package proof

import (
	"execution-service/circuit"
	"execution-service/instance"
	"execution-service/message"
	"execution-service/model"
	"execution-service/parameters"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type ProofService struct {
	proofParameters     parameters.ProofParameters
	SignatureParameters parameters.SignatureParameters
	InstanceService     instance.InstanceService
	ModelService        model.ModelService
	MessageService      message.MessageService
}

func InitializeProofService() ProofService {
	proofParameters := parameters.NewProofParameters()
	signatureParameters := parameters.NewSignatureParameters()
	instanceService := instance.NewInstanceService()
	modelService := model.NewModelService()
	messageService := message.NewMessageService()
	return NewProofService(proofParameters, signatureParameters, instanceService, modelService, messageService)
}

func NewProofService(proofParameters parameters.ProofParameters, signatureParameters parameters.SignatureParameters, instanceService instance.InstanceService, modelService model.ModelService, messageService message.MessageService) ProofService {
	fmt.Printf("Instantiation constraint system has %d constraints\n", proofParameters.CsInstantiation.GetNbConstraints())
	fmt.Printf("Transition constraint system has %d constraints\n", proofParameters.CsTransition.GetNbConstraints())
	fmt.Printf("Termination constraint system has %d constraints\n", proofParameters.CsTermination.GetNbConstraints())
	return ProofService{
		proofParameters:     proofParameters,
		SignatureParameters: signatureParameters,
		InstanceService:     instanceService,
		ModelService:        modelService,
		MessageService:      messageService,
	}
}

func (service *ProofService) ProveInstantiation(cmd ProveInstantiationCommand) (Proof, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return Proof{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	signature := instance.Sign(privateKey)
	assignment := &circuit.InstantiationCircuit{
		Model:          circuit.FromModel(model),
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	groth16Proof, err := groth16.Prove(service.proofParameters.CsInstantiation, service.proofParameters.PkInstantiation, witness)
	if err != nil {
		return Proof{}, err
	}
	return toProof(groth16Proof, model.Hash, instance.Hash)
}

func (service *ProofService) ProveTransition(cmd ProveTransitionCommand) (Proof, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	currentInstance, err := service.InstanceService.FindInstanceById(cmd.CurrentInstance)
	if err != nil {
		return Proof{}, err
	}
	nextInstance, err := service.InstanceService.FindInstanceById(cmd.NextInstance)
	if err != nil {
		return Proof{}, err
	}
	transition, err := model.FindTransitionById(cmd.Transition)
	if err != nil {
		return Proof{}, err
	}
	constraintInput, err := service.MessageService.FindConstraintInput(transition.Constraint, currentInstance)
	if err != nil {
		return Proof{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	nextSignature := nextInstance.Sign(privateKey)
	assignment := &circuit.TransitionCircuit{
		Model:           circuit.FromModel(model),
		CurrentInstance: circuit.FromInstance(currentInstance),
		NextInstance:    circuit.FromInstance(nextInstance),
		Transition:      circuit.ToTransition(model, transition),
		Authentication:  circuit.ToAuthentication(nextInstance, nextSignature),
		ConstraintInput: circuit.FromConstraintInput(constraintInput),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTransition, service.proofParameters.PkTransition, witness)
	if err != nil {
		return Proof{}, err
	}

	return toProof(proof, model.Hash, currentInstance.Hash, nextInstance.Hash)
}

func (service *ProofService) ProveTermination(cmd ProveTerminationCommand) (Proof, error) {
	model, err := service.ModelService.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return Proof{}, err
	}
	privateKey := service.SignatureParameters.GetPrivateKeyForIdentity(cmd.Identity)
	signature := instance.Sign(privateKey)
	assignment := &circuit.TerminationCircuit{
		Model:          circuit.FromModel(model),
		Instance:       circuit.FromInstance(instance),
		Authentication: circuit.ToAuthentication(instance, signature),
		EndPlaceProof:  circuit.ToEndPlaceProof(model, cmd.EndPlace),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTermination, service.proofParameters.PkTermination, witness)
	if err != nil {
		return Proof{}, err
	}
	return toProof(proof, model.Hash, instance.Hash)
}
