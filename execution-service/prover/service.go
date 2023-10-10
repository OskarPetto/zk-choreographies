package prover

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

type ProverService struct {
	proofParameters     parameters.ProverParameters
	SignatureParameters parameters.SignatureParameters
	InstanceService     instance.InstanceService
	ModelService        model.ModelService
	MessageService      message.MessageService
}

func InitializeProverService() ProverService {
	proofParameters := parameters.NewProverParameters()
	signatureParameters := parameters.NewSignatureParameters()
	modelService := model.NewModelService()
	messageService := message.NewMessageService()
	instanceService := instance.NewInstanceService(modelService, messageService)
	return NewProverService(proofParameters, signatureParameters, instanceService)
}

func NewProverService(proofParameters parameters.ProverParameters, signatureParameters parameters.SignatureParameters, instanceService instance.InstanceService) ProverService {
	fmt.Printf("Instantiation constraint system has %d constraints\n", proofParameters.CsInstantiation.GetNbConstraints())
	fmt.Printf("Transition constraint system has %d constraints\n", proofParameters.CsTransition.GetNbConstraints())
	fmt.Printf("Termination constraint system has %d constraints\n", proofParameters.CsTermination.GetNbConstraints())
	return ProverService{
		proofParameters:     proofParameters,
		SignatureParameters: signatureParameters,
		InstanceService:     instanceService,
		ModelService:        instanceService.ModelService,
		MessageService:      instanceService.MessageService,
	}
}

func (service *ProverService) ProveInstantiation(cmd ProveInstantiationCommand) (Proof, error) {
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

func (service *ProverService) ProveTransition(cmd ProveTransitionCommand) (Proof, error) {
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

func (service *ProverService) ProveTermination(cmd ProveTerminationCommand) (Proof, error) {
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
