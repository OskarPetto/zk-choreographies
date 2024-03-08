package prover

import (
	"execution-service/circuit"
	"execution-service/parameters"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

type IProverService interface {
	ProveInstantiation(cmd ProveInstantiationCommand) (Proof, error)
	ProveTransition(cmd ProveTransitionCommand) (Proof, error)
	ProveTermination(cmd ProveTerminationCommand) (Proof, error)
}

type ProverService struct {
	proofParameters parameters.ProverParameters
}

func InitializeProverService() ProverService {
	proofParameters := parameters.NewProverParameters()
	return NewProverService(proofParameters)
}

func NewProverService(proofParameters parameters.ProverParameters) ProverService {
	fmt.Printf("Instantiation constraint system has %d constraints\n", proofParameters.CsInstantiation.GetNbConstraints())
	fmt.Printf("Transition constraint system has %d constraints\n", proofParameters.CsTransition.GetNbConstraints())
	fmt.Printf("Termination constraint system has %d constraints\n", proofParameters.CsTermination.GetNbConstraints())
	return ProverService{
		proofParameters: proofParameters,
	}
}

func (service ProverService) ProveInstantiation(cmd ProveInstantiationCommand) (Proof, error) {
	assignment := &circuit.InstantiationCircuit{
		Model:          circuit.FromModel(cmd.Model),
		Instance:       circuit.FromInstance(cmd.Instance),
		Authentication: circuit.ToAuthentication(cmd.Instance, cmd.Signature),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	groth16Proof, err := groth16.Prove(service.proofParameters.CsInstantiation, service.proofParameters.PkInstantiation, witness)
	if err != nil {
		return Proof{}, err
	}
	return toProof(groth16Proof, cmd.Instance.SaltedHash.Hash)
}

func (service ProverService) ProveTransition(cmd ProveTransitionCommand) (Proof, error) {

	senderAuthentication := circuit.ToAuthentication(cmd.NextInstance, cmd.InitiatingParticipantSignature)
	recipientAuthentication := senderAuthentication
	if cmd.RespondingParticipantSignature != nil {
		recipientAuthentication = circuit.ToAuthentication(cmd.NextInstance, *cmd.RespondingParticipantSignature)
	}

	assignment := &circuit.TransitionCircuit{
		Model:                               circuit.FromModel(cmd.Model),
		CurrentInstance:                     circuit.FromInstance(cmd.CurrentInstance),
		NextInstance:                        circuit.FromInstance(cmd.NextInstance),
		Transition:                          circuit.ToTransition(cmd.Model, cmd.Transition),
		InitiatingParticipantAuthentication: senderAuthentication,
		RespondingParticipantAuthentication: recipientAuthentication,
		ConditionInput:                      circuit.FromConditionInput(cmd.ConditionInput),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTransition, service.proofParameters.PkTransition, witness)
	if err != nil {
		return Proof{}, err
	}

	return toProof(proof, cmd.CurrentInstance.SaltedHash.Hash, cmd.NextInstance.SaltedHash.Hash)
}

func (service ProverService) ProveTermination(cmd ProveTerminationCommand) (Proof, error) {

	assignment := &circuit.TerminationCircuit{
		Model:          circuit.FromModel(cmd.Model),
		Instance:       circuit.FromInstance(cmd.Instance),
		Authentication: circuit.ToAuthentication(cmd.Instance, cmd.Signature),
		EndPlaceProof:  circuit.ToEndPlaceProof(cmd.Model, cmd.Instance),
	}
	witness, err := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
	if err != nil {
		return Proof{}, err
	}
	proof, err := groth16.Prove(service.proofParameters.CsTermination, service.proofParameters.PkTermination, witness)
	if err != nil {
		return Proof{}, err
	}
	return toProof(proof, cmd.Instance.SaltedHash.Hash)
}
