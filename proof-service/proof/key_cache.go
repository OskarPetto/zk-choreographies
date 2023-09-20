package proof

import (
	"fmt"
	"io"
	"os"
	"proof-service/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

const csPathInstantiation = "./files/instantiation.constraint_system"
const csPathTransition = "./files/transition.constraint_system"
const csPathTermination = "./files/termination.constraint_system"
const pkPathInstantiation = "./files/instantiation.proving_key"
const pkPathTransition = "./files/transition.proving_key"
const pkPathTermination = "./files/termination.proving_key"

type KeyCache struct {
	csInstantiation constraint.ConstraintSystem
	csTransition    constraint.ConstraintSystem
	csTermination   constraint.ConstraintSystem
	pkInstantiation groth16.ProvingKey
	pkTransition    groth16.ProvingKey
	pkTermination   groth16.ProvingKey
}

func NewKeyCache() KeyCache {
	csInstantiation := importConstraintSystem(&circuit.InstantiationCircuit{}, csPathInstantiation)
	csTransition := importConstraintSystem(&circuit.TransitionCircuit{}, csPathTransition)
	csTermination := importConstraintSystem(&circuit.TerminationCircuit{}, csPathTermination)
	pkInstantiation := importProvingKey(csInstantiation, pkPathInstantiation)
	pkTransition := importProvingKey(csTransition, pkPathTransition)
	pkTermination := importProvingKey(csTermination, pkPathTermination)
	return KeyCache{
		csInstantiation,
		csTransition,
		csTermination,
		pkInstantiation,
		pkTransition,
		pkTermination,
	}
}

func importConstraintSystem(circuit frontend.Circuit, path string) constraint.ConstraintSystem {
	cs := groth16.NewCS(ecc.BN254)
	err := readFile(cs, path)
	if err != nil {
		cs = compileCircuit(circuit, path)
	}
	return cs
}

func importProvingKey(cs constraint.ConstraintSystem, path string) groth16.ProvingKey {
	pk := groth16.NewProvingKey(ecc.BN254)
	err := readFile(pk, path)
	if err != nil {
		pk = generateProvingKey(cs, path)
	}
	return pk
}

func compileCircuit(circuit frontend.Circuit, path string) constraint.ConstraintSystem {
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	check(err)
	writeFile(cs, path)
	return cs
}

func generateProvingKey(cs constraint.ConstraintSystem, path string) groth16.ProvingKey {
	pk, _, err := groth16.Setup(cs)
	check(err)
	writeFile(pk, path)
	return pk
}

func writeFile(writeable io.WriterTo, path string) {
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	bytesWritten, err := writeable.WriteTo(file)
	check(err)
	fmt.Printf("Wrote file of size %d in %s\n", bytesWritten, path)
}

func readFile(readable io.ReaderFrom, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	bytesRead, err := readable.ReadFrom(file)
	fmt.Printf("Read file of size %d in %s\n", bytesRead, path)
	return err
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
