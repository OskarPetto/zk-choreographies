package proof

import (
	"execution-service/authentication"
	"execution-service/domain"
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

func (service *ProofService) ProveInstantiation(model domain.Model, instance domain.Instance, signature authentication.Signature) (Proof, error) {
	assignment := &circuit.InstantiationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
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

func (service *ProofService) ProveTransition(model domain.Model, currentInstance domain.Instance, nextInstance domain.Instance, nextSignature authentication.Signature) (Proof, error) {
	assignment := &circuit.TransitionCircuit{
		CurrentInstance:       circuit.FromInstance(currentInstance),
		NextInstance:          circuit.FromInstance(nextInstance),
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuit.FromModel(model),
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

func (service *ProofService) ProveTermination(model domain.Model, instance domain.Instance, signature authentication.Signature) (Proof, error) {
	assignment := &circuit.TerminationCircuit{
		Instance:  circuit.FromInstance(instance),
		Signature: circuit.FromSignature(signature),
		Model:     circuit.FromModel(model),
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
