package proof

import (
	"bytes"
	"proof-service/circuit"
	"proof-service/commitment"
	"proof-service/workflow"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type ProofService struct {
	keyCache KeyCache
}

func NewProofService() ProofService {
	return ProofService{
		keyCache: NewKeyCache(),
	}
}

func (service *ProofService) ProveInstantiation(instance workflow.Instance, comm commitment.Commitment, pertiNet workflow.PetriNet) ([]byte, error) {
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
		Commitment: circuit.FromCommitment(comm),
		PetriNet:   circuitPetriNet,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return []byte{}, err
	}
	proof, err := groth16.Prove(service.keyCache.csInstantiation, service.keyCache.pkInstantiation, witness)
	if err != nil {
		return []byte{}, err
	}
	byteBuffer := new(bytes.Buffer)
	proof.WriteTo(byteBuffer)
	return byteBuffer.Bytes(), nil
}
