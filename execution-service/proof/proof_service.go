package proof

import (
	"bytes"
	"math/big"
	"proof-service/authentication"
	"proof-service/circuit"
	"proof-service/domain"
	"proof-service/proof/parameters"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

type VerifierInput struct {
	A     [2]*big.Int
	B     [2][2]*big.Int
	C     [2]*big.Int
	Input []*big.Int
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

func (service *ProofService) ProveInstantiation(model domain.Model, instance domain.Instance, signature authentication.Signature) (VerifierInput, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return VerifierInput{}, err
	}
	circuitModel, err := circuit.FromModel(model)
	if err != nil {
		return VerifierInput{}, err
	}
	assignment := &circuit.InstantiationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     circuitModel,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return VerifierInput{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsInstantiation, service.proofParameters.PkInstantiation, witness)
	if err != nil {
		return VerifierInput{}, err
	}
	return toVerifierInput(proof, witness)
}

func (service *ProofService) ProveTransition(model domain.Model, currentInstance domain.Instance, nextInstance domain.Instance, nextSignature authentication.Signature) (VerifierInput, error) {
	currentCircuitInstance, err := circuit.FromInstance(currentInstance)
	if err != nil {
		return VerifierInput{}, err
	}
	nextCircuitInstance, err := circuit.FromInstance(nextInstance)
	if err != nil {
		return VerifierInput{}, err
	}
	circuitModel, err := circuit.FromModel(model)
	if err != nil {
		return VerifierInput{}, err
	}

	assignment := &circuit.TransitionCircuit{
		CurrentInstance:       currentCircuitInstance,
		NextInstance:          nextCircuitInstance,
		NextInstanceSignature: circuit.FromSignature(nextSignature),
		Model:                 circuitModel,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return VerifierInput{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTransition, service.proofParameters.PkTransition, witness)
	if err != nil {
		return VerifierInput{}, err
	}
	return toVerifierInput(proof, witness)
}

func (service *ProofService) ProveTermination(model domain.Model, instance domain.Instance, signature authentication.Signature) (VerifierInput, error) {
	circuitInstance, err := circuit.FromInstance(instance)
	if err != nil {
		return VerifierInput{}, err
	}
	circuitModel, err := circuit.FromModel(model)
	if err != nil {
		return VerifierInput{}, err
	}
	assignment := &circuit.TerminationCircuit{
		Instance:  circuitInstance,
		Signature: circuit.FromSignature(signature),
		Model:     circuitModel,
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return VerifierInput{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTermination, service.proofParameters.PkTermination, witness)
	if err != nil {
		return VerifierInput{}, err
	}
	return toVerifierInput(proof, witness)
}

func toVerifierInput(proof groth16.Proof, fullWitness witness.Witness) (VerifierInput, error) {
	var buf bytes.Buffer
	proof.WriteRawTo(&buf)
	proofBytes := buf.Bytes()

	var verifierInput VerifierInput

	// proof.Ar, proof.Bs, proof.Krs
	const fpSize = 4 * 8
	verifierInput.A[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	verifierInput.A[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	verifierInput.B[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	verifierInput.B[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	verifierInput.B[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	verifierInput.B[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	verifierInput.C[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	verifierInput.C[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	publicWitness, err := fullWitness.Public()
	if err != nil {
		return VerifierInput{}, err
	}
	publicWitnessBytes, err := publicWitness.MarshalBinary()
	if err != nil {
		return VerifierInput{}, err
	}

	fieldElementCount := len(publicWitnessBytes) / fr.Bytes
	verifierInput.Input = make([]*big.Int, fieldElementCount)
	for i := 0; i < fieldElementCount; i++ {
		var fieldElement [fr.Bytes]byte
		copy(fieldElement[:], publicWitnessBytes[i*fr.Bytes:(i+1)*fr.Bytes-1])
		verifierInput.Input[i] = new(big.Int).SetBytes(fieldElement[:])
	}

	return verifierInput, nil
}
