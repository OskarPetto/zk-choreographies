package proof

import (
	"bytes"
	"proof-service/circuit"
	"proof-service/crypto"
	"proof-service/domain"
	"proof-service/proof/parameters"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type ProofService struct {
	isLoaded         bool
	proofParameters  parameters.ProofParameters
	signatureService crypto.SignatureService
}

var proofService ProofService

func NewProofService() ProofService {
	if !proofService.isLoaded {
		proofService = ProofService{
			isLoaded:         true,
			proofParameters:  parameters.LoadProofParameters(),
			signatureService: crypto.NewSignatureService(),
		}
	}
	return proofService
}

func (service *ProofService) ProveInstantiation(instance domain.Instance, pertiNet domain.PetriNet) ([]byte, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return []byte{}, err
	}
	circuitPetriNet, err := circuit.FromPetriNet(pertiNet)
	if err != nil {
		return []byte{}, err
	}
	saltedHash := crypto.HashInstance(instance)
	assignment := &circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(saltedHash),
		Signature:  circuit.FromSignature(service.signatureService.Sign(saltedHash)),
		PetriNet:   circuitPetriNet,
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

func (service *ProofService) ProveTransition(currentInstance domain.Instance, nextInstance domain.Instance, pertiNet domain.PetriNet) ([]byte, error) {
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
	currentSaltedHash := crypto.HashInstance(currentInstance)
	nextSaltedHash := crypto.HashInstance(nextInstance)

	assignment := &circuit.TransitionCircuit{
		CurrentInstance:           currentCircuitInstance,
		CurrentInstanceSaltedHash: circuit.FromSaltedHash(currentSaltedHash),
		NextInstance:              nextCircuitInstance,
		NextInstanceSaltedHash:    circuit.FromSaltedHash(nextSaltedHash),
		NextInstanceSignature:     circuit.FromSignature(service.signatureService.Sign(nextSaltedHash)),
		PetriNet:                  circuitPetriNet,
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

func (service *ProofService) ProveTermination(instance domain.Instance, pertiNet domain.PetriNet) ([]byte, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return []byte{}, err
	}
	circuitPetriNet, err := circuit.FromPetriNet(pertiNet)
	if err != nil {
		return []byte{}, err
	}
	saltedHash := crypto.HashInstance(instance)
	assignment := &circuit.TerminationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(saltedHash),
		Signature:  circuit.FromSignature(service.signatureService.Sign(saltedHash)),
		PetriNet:   circuitPetriNet,
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
