//go:build wireinject
// +build wireinject

package proof

import (
	"execution-service/instance"
	"execution-service/model"
	"execution-service/proof/circuit"
	"execution-service/signature"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/google/wire"
)

type ProofService struct {
	proofParameters  ProofParameters
	InstanceService  instance.InstanceService
	ModelService     model.ModelService
	signatureService signature.SignatureService
}

func InitializeProofService(modelPort model.ModelPort) ProofService {
	wire.Build(instance.NewInstanceService, model.NewModelService, signature.InitializeSignatureService, NewProofParameters, NewProofService)
	return ProofService{}
}

func NewProofService(proofParameters ProofParameters, instanceService instance.InstanceService, modelService model.ModelService, signatureService signature.SignatureService) ProofService {
	return ProofService{
		proofParameters:  proofParameters,
		InstanceService:  instanceService,
		ModelService:     modelService,
		signatureService: signatureService,
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
	signature := service.signatureService.Sign(instance)
	assignment := &circuit.InstantiationCircuit{
		Model:     circuit.FromModel(model),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsInstantiation, service.proofParameters.PkInstantiation, witness)
	if err != nil {
		return Proof{}, err
	}
	return newProof(proof, witness)
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
	signature := service.signatureService.Sign(nextInstance)
	assignment := &circuit.TransitionCircuit{
		Model:                 circuit.FromModel(model),
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(signature),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTransition, service.proofParameters.PkTransition, witness)
	if err != nil {
		return Proof{}, err
	}
	return newProof(proof, witness)
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
	signature := service.signatureService.Sign(instance)
	assignment := &circuit.TerminationCircuit{
		Model:     circuit.FromModel(model),
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTermination, service.proofParameters.PkTermination, witness)
	if err != nil {
		return Proof{}, err
	}
	return newProof(proof, witness)
}
