package proof

import (
	"bytes"
	"proof-service/authentication"
	"proof-service/circuit"
	"proof-service/domain"
	"proof-service/proof/parameters"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type Proof struct {
	Value []byte
}

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

func (service *ProofService) ProveInstantiation(instance domain.Instance, pertiNet domain.Model, signature authentication.Signature) (Proof, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return Proof{}, err
	}
	circuitModel, err := circuit.FromModel(pertiNet)
	if err != nil {
		return Proof{}, err
	}
	assignment := &circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     circuitModel,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsInstantiation, service.proofParameters.PkInstantiation, witness)
	if err != nil {
		return Proof{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return Proof{
		Value: byteBuffer.Bytes(),
	}, nil
}

func (service *ProofService) ProveTransition(currentInstance domain.Instance, nextInstance domain.Instance, pertiNet domain.Model, nextSignature authentication.Signature) (Proof, error) {
	currentCircuitInstance, err := circuit.FromInstance(currentInstance)
	if err != nil {
		return Proof{}, err
	}
	nextCircuitInstance, err := circuit.FromInstance(nextInstance)
	if err != nil {
		return Proof{}, err
	}
	circuitModel, err := circuit.FromModel(pertiNet)
	if err != nil {
		return Proof{}, err
	}

	assignment := &circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuitModel,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTransition, service.proofParameters.PkTransition, witness)
	if err != nil {
		return Proof{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return Proof{
		Value: byteBuffer.Bytes(),
	}, nil
}

func (service *ProofService) ProveTermination(instance domain.Instance, pertiNet domain.Model, signature authentication.Signature) (Proof, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return Proof{}, err
	}
	circuitModel, err := circuit.FromModel(pertiNet)
	if err != nil {
		return Proof{}, err
	}
	assignment := &circuit.TerminationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     circuitModel,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTermination, service.proofParameters.PkTermination, witness)
	if err != nil {
		return Proof{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return Proof{
		Value: byteBuffer.Bytes(),
	}, nil
}
