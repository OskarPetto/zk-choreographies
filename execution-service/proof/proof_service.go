//go:build wireinject
// +build wireinject

package proof

import (
	"execution-service/hash"
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
	modelPort        model.ModelPort
	InstanceService  instance.InstanceService
	HashService      hash.HashService
	signatureService signature.SignatureService
}

func InitializeProofService(modelPort model.ModelPort) ProofService {
	wire.Build(instance.NewInstanceService, hash.NewHashService, signature.InitializeSignatureService, NewProofParameters, NewProofService)
	return ProofService{}
}

func NewProofService(proofParameters ProofParameters, modelPort model.ModelPort, instanceService instance.InstanceService, hashService hash.HashService, signatureService signature.SignatureService) ProofService {
	return ProofService{
		proofParameters:  proofParameters,
		modelPort:        modelPort,
		InstanceService:  instanceService,
		HashService:      hashService,
		signatureService: signatureService,
	}
}

func (service *ProofService) ProveInstantiation(cmd ProveInstantiationCommand) (Proof, error) {
	model, err := service.modelPort.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	modelHash, err := service.HashService.FindHashByModelId(model.Id)
	if err != nil {
		return Proof{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return Proof{}, err
	}
	signature := service.signatureService.Sign(instance)
	assignment := &circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(modelHash),
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
	model, err := service.modelPort.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	modelHash, err := service.HashService.FindHashByModelId(model.Id)
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
		ModelHash:             circuit.FromHash(modelHash),
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
	model, err := service.modelPort.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	modelHash, err := service.HashService.FindHashByModelId(model.Id)
	if err != nil {
		return Proof{}, err
	}
	instance, err := service.InstanceService.FindInstanceById(cmd.Instance)
	if err != nil {
		return Proof{}, err
	}
	signature := service.signatureService.Sign(instance)
	assignment := &circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(modelHash),
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
