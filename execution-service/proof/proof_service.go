package proof

import (
	"execution-service/proof/circuit"
	"execution-service/proof/parameters"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type ProofService struct {
	isLoaded        bool
	proofParameters parameters.ProofParameters
}

var proofService ProofService

func NewProofService() ProofService {
	if !proofService.isLoaded {
		proofService = ProofService{
			isLoaded:        true,
			proofParameters: parameters.LoadProofParameters(),
		}
	}
	return proofService
}

func (service *ProofService) ProveInstantiation(cmd ProveInstantiationCommand) (Proof, error) {
	assignment := &circuit.InstantiationCircuit{
		ModelHash: circuit.FromHash(cmd.ModelHash),
		Model:     circuit.FromModel(cmd.Model),
		Instance:  circuit.FromInstance(cmd.Instance),
		Signature: circuit.FromSignature(cmd.Signature),
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
	assignment := &circuit.TransitionCircuit{
		ModelHash:             circuit.FromHash(cmd.ModelHash),
		Model:                 circuit.FromModel(cmd.Model),
		CurrentInstance:       circuit.FromInstance(cmd.CurrentInstance),
		NextInstance:          circuit.FromInstance(cmd.NextInstance),
		NextInstanceSignature: circuit.FromSignature(cmd.NextSignature),
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
	assignment := &circuit.TerminationCircuit{
		ModelHash: circuit.FromHash(cmd.ModelHash),
		Model:     circuit.FromModel(cmd.Model),
		Instance:  circuit.FromInstance(cmd.Instance),
		Signature: circuit.FromSignature(cmd.Signature),
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
