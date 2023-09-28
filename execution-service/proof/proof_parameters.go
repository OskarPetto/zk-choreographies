package proof

import (
	"bytes"
	"execution-service/file"
	"execution-service/proof/circuit"
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

const instantiationCsFilename = "instantiation.constraint_system"
const transitionCsFilename = "transition.constraint_system"
const terminationCsFilename = "termination.constraint_system"
const instantiationPkFilename = "instantiation.proving_key"
const transitionPkFilename = "transition.proving_key"
const terminationPkFilename = "termination.proving_key"
const instantiationVkFilename = "instantiation.sol"
const transitionVkFilename = "transition.sol"
const terminationVkFilename = "termination.sol"

type ProofParameters struct {
	CsInstantiation constraint.ConstraintSystem
	CsTransition    constraint.ConstraintSystem
	CsTermination   constraint.ConstraintSystem
	PkInstantiation groth16.ProvingKey
	PkTransition    groth16.ProvingKey
	PkTermination   groth16.ProvingKey
}

func NewProofParameters() ProofParameters {
	csInstantiation := importConstraintSystem(&circuit.InstantiationCircuit{}, instantiationCsFilename)
	csTransition := importConstraintSystem(&circuit.TransitionCircuit{}, transitionCsFilename)
	csTermination := importConstraintSystem(&circuit.TerminationCircuit{}, terminationCsFilename)
	pkInstantiation := importProvingKey(csInstantiation, instantiationPkFilename, instantiationVkFilename)
	pkTransition := importProvingKey(csTransition, transitionPkFilename, transitionVkFilename)
	pkTermination := importProvingKey(csTermination, terminationPkFilename, terminationVkFilename)
	return ProofParameters{
		csInstantiation,
		csTransition,
		csTermination,
		pkInstantiation,
		pkTransition,
		pkTermination,
	}
}

func importConstraintSystem(circuit frontend.Circuit, filename string) constraint.ConstraintSystem {
	cs := groth16.NewCS(ecc.BN254)
	err := file.ReadPublicFile(cs, filename)
	if err != nil {
		cs = compileCircuit(circuit, filename)
	}
	return cs
}

func importProvingKey(cs constraint.ConstraintSystem, pkFilename string, vkFilename string) groth16.ProvingKey {
	pk := groth16.NewProvingKey(ecc.BN254)
	err := file.ReadPublicFile(pk, pkFilename)
	if err != nil {
		pk = generateProvingKey(cs, pkFilename, vkFilename)
	}
	return pk
}

func compileCircuit(circuit frontend.Circuit, filename string) constraint.ConstraintSystem {
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	utils.PanicOnError(err)
	file.WritePublicFile(cs, filename)
	return cs
}

func generateProvingKey(cs constraint.ConstraintSystem, pkFilename string, vkFilename string) groth16.ProvingKey {
	pk, vk, err := groth16.Setup(cs)
	utils.PanicOnError(err)
	file.WritePublicFile(pk, pkFilename)
	byteBuffer := new(bytes.Buffer)
	vk.ExportSolidity(byteBuffer)
	file.WritePublicFile(byteBuffer, vkFilename)
	return pk
}
