package proof

import (
	"proof-service/circuit"
	"proof-service/file"

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

type Parameters struct {
	isLoaded        bool
	csInstantiation constraint.ConstraintSystem
	csTransition    constraint.ConstraintSystem
	csTermination   constraint.ConstraintSystem
	pkInstantiation groth16.ProvingKey
	pkTransition    groth16.ProvingKey
	pkTermination   groth16.ProvingKey
}

var parameters Parameters

func LoadParameters() Parameters {
	if !parameters.isLoaded {
		csInstantiation := importConstraintSystem(&circuit.InstantiationCircuit{}, instantiationCsFilename)
		csTransition := importConstraintSystem(&circuit.TransitionCircuit{}, transitionCsFilename)
		csTermination := importConstraintSystem(&circuit.TerminationCircuit{}, terminationCsFilename)
		pkInstantiation := importProvingKey(csInstantiation, instantiationPkFilename)
		pkTransition := importProvingKey(csTransition, transitionPkFilename)
		pkTermination := importProvingKey(csTermination, terminationPkFilename)
		parameters = Parameters{
			true,
			csInstantiation,
			csTransition,
			csTermination,
			pkInstantiation,
			pkTransition,
			pkTermination,
		}
	}
	return parameters
}

func importConstraintSystem(circuit frontend.Circuit, filename string) constraint.ConstraintSystem {
	cs := groth16.NewCS(ecc.BN254)
	err := file.ReadFile(cs, filename)
	if err != nil {
		cs = compileCircuit(circuit, filename)
	}
	return cs
}

func importProvingKey(cs constraint.ConstraintSystem, filename string) groth16.ProvingKey {
	pk := groth16.NewProvingKey(ecc.BN254)
	err := file.ReadFile(pk, filename)
	if err != nil {
		pk = generateProvingKey(cs, filename)
	}
	return pk
}

func compileCircuit(circuit frontend.Circuit, filename string) constraint.ConstraintSystem {
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	check(err)
	file.WriteFile(cs, filename)
	return cs
}

func generateProvingKey(cs constraint.ConstraintSystem, filename string) groth16.ProvingKey {
	pk, _, err := groth16.Setup(cs)
	check(err)
	file.WriteFile(pk, filename)
	return pk
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
