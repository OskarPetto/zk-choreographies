package proof

import (
	"bytes"
	"proof-service/circuit"
	"proof-service/crypto"
	"proof-service/workflow"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

func ProveInstantiation(instance workflow.Instance, pertiNet workflow.PetriNet) ([]byte, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return []byte{}, err
	}
	circuitPetriNet, err := circuit.FromPetriNet(pertiNet)
	if err != nil {
		return []byte{}, err
	}
	assignment := &circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		Commitment: circuit.FromCommitment(crypto.Commit(instance)),
		PetriNet:   circuitPetriNet,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return []byte{}, err
	}
	proof, err := groth16.Prove(LoadProofParameters().csInstantiation, LoadProofParameters().pkInstantiation, witness)
	if err != nil {
		return []byte{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return byteBuffer.Bytes(), nil
}

func ProveTransition(currentInstance workflow.Instance, nextInstance workflow.Instance, pertiNet workflow.PetriNet) ([]byte, error) {
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
		CurrentInstance:   currentCircuitInstance,
		CurrentCommitment: circuit.FromCommitment(crypto.Commit(currentInstance)),
		NextInstance:      nextCircuitInstance,
		NextCommitment:    circuit.FromCommitment(crypto.Commit(nextInstance)),
		PetriNet:          circuitPetriNet,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return []byte{}, err
	}
	proof, err := groth16.Prove(LoadProofParameters().csTransition, LoadProofParameters().pkTransition, witness)
	if err != nil {
		return []byte{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return byteBuffer.Bytes(), nil
}

func ProveTermination(instance workflow.Instance, pertiNet workflow.PetriNet) ([]byte, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return []byte{}, err
	}
	circuitPetriNet, err := circuit.FromPetriNet(pertiNet)
	if err != nil {
		return []byte{}, err
	}
	assignment := &circuit.TerminationCircuit{
		Instance:   circuitInstance,
		Commitment: circuit.FromCommitment(crypto.Commit(instance)),
		PetriNet:   circuitPetriNet,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return []byte{}, err
	}
	proof, err := groth16.Prove(LoadProofParameters().csTermination, LoadProofParameters().pkTermination, witness)
	if err != nil {
		return []byte{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return byteBuffer.Bytes(), nil
}
