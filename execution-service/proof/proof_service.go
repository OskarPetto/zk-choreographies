package proof

import (
	"execution-service/authentication"
	"execution-service/circuit"
	"execution-service/domain"
	"execution-service/proof/parameters"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type ProofService struct {
	isLoaded         bool
	proofParameters  parameters.ProofParameters
	hashService      domain.HashService
	instanceService  domain.InstanceService
	signatureService authentication.SignatureService
	modelService     domain.ModelService
}

var proofService ProofService

func NewProofService() ProofService {
	if !proofService.isLoaded {
		proofService = ProofService{
			isLoaded:         true,
			proofParameters:  parameters.LoadProofParameters(),
			hashService:      domain.NewHashService(),
			instanceService:  domain.NewInstanceService(),
			signatureService: authentication.NewSignatureService(),
			modelService:     domain.ModelServiceImpl,
		}
	}
	return proofService
}

func (service *ProofService) ProveInstantiation(cmd ProveInstantiationCommand) (Proof, error) {
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	modelHash, err := service.hashService.FindHashByModelId(model.Id)
	if err != nil {
		return Proof{}, err
	}
	instance, err := service.instanceService.FindInstanceById(cmd.Instance)
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
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	modelHash, err := service.hashService.FindHashByModelId(model.Id)
	if err != nil {
		return Proof{}, err
	}
	currentInstance, err := service.instanceService.FindInstanceById(cmd.CurrentInstance)
	if err != nil {
		return Proof{}, err
	}
	nextInstance, err := service.instanceService.FindInstanceById(cmd.NextInstance)
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
	model, err := service.modelService.FindModelById(cmd.Model)
	if err != nil {
		return Proof{}, err
	}
	modelHash, err := service.hashService.FindHashByModelId(model.Id)
	if err != nil {
		return Proof{}, err
	}
	instance, err := service.instanceService.FindInstanceById(cmd.Instance)
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
