package proof

import (
	"bytes"
	"proof-service/authentication"
	"proof-service/circuit"
	"proof-service/instance"
	"proof-service/model"
	"proof-service/proof/parameters"

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

func (service *ProofService) ProveInstantiation(instance instance.Instance, pertiNet model.PetriNet, signature authentication.Signature) ([]byte, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return []byte{}, err
	}
	circuitPetriNet, err := circuit.FromPetriNet(pertiNet)
	if err != nil {
		return []byte{}, err
	}
	assignment := &circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		PetriNet:  circuitPetriNet,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return []byte{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsInstantiation, service.proofParameters.PkInstantiation, witness)
	if err != nil {
		return []byte{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return byteBuffer.Bytes(), nil
}

func (service *ProofService) ProveTransition(currentInstance instance.Instance, nextInstance instance.Instance, pertiNet model.PetriNet, nextSignature authentication.Signature) ([]byte, error) {
	currentCircuitInstance, err := circuit.FromInstance(currentInstance)
	if err != nil {
		return []byte{}, err
	}
	nextCircuitInstance, err := circuit.FromInstance(nextInstance)
	if err != nil {
		return []byte{}, err
	}
	circuitPetriNet, err := circuit.FromPetriNet(pertiNet)
	if err != nil {
		return []byte{}, err
	}

	assignment := &circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		PetriNet:              circuitPetriNet,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return []byte{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTransition, service.proofParameters.PkTransition, witness)
	if err != nil {
		return []byte{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return byteBuffer.Bytes(), nil
}

func (service *ProofService) ProveTermination(instance instance.Instance, pertiNet model.PetriNet, signature authentication.Signature) ([]byte, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return []byte{}, err
	}
	circuitPetriNet, err := circuit.FromPetriNet(pertiNet)
	if err != nil {
		return []byte{}, err
	}
	assignment := &circuit.TerminationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		PetriNet:  circuitPetriNet,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return []byte{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTermination, service.proofParameters.PkTermination, witness)
	if err != nil {
		return []byte{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return byteBuffer.Bytes(), nil
}
