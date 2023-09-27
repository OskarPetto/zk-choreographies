package proof

import (
	"bytes"
	"math/big"
	"proof-service/authentication"
	"proof-service/domain"
	"proof-service/proof/circuit"
	"proof-service/proof/parameters"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

type Proof struct {
	A           [2]*big.Int
	B           [2][2]*big.Int
	C           [2]*big.Int
	PublicInput []*big.Int
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
	return toProof(proof, witness)
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
	return toProof(proof, witness)
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
	return toProof(proof, witness)
}

func toProof(groth16Proof groth16.Proof, fullWitness witness.Witness) (Proof, error) {

	var proof Proof

	var buf bytes.Buffer
	groth16Proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()

	// proof.Ar, proof.Bs, proof.Krs
	const fpSize = 4 * 8
	proof.A[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	proof.A[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	proof.B[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	proof.B[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	proof.B[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	proof.B[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	proof.C[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	proof.C[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	publicWitness, err := fullWitness.Public()
	if err != nil {
		return Proof{}, err
	}
	publicWitnessBytes, err := publicWitness.MarshalBinary()
	if err != nil {
		return Proof{}, err
	}

	headerByteCount := 12
	fieldElementCount := (len(publicWitnessBytes) - headerByteCount) / fr.Bytes
	proof.PublicInput = make([]*big.Int, fieldElementCount)
	for i := 0; i < fieldElementCount; i++ {
		var fieldElement [fr.Bytes]byte
		copy(fieldElement[:], publicWitnessBytes[headerByteCount+i*fr.Bytes:headerByteCount+(i+1)*fr.Bytes-1])
		proof.PublicInput[i] = new(big.Int).SetBytes(fieldElement[:])
	}

	return proof, nil
}
