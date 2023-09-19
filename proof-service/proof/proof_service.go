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

const instantiationPkFilename = "instantiation.proving_key"
const instantiationCsFilename = "instantiation.constraint_system"
const transitionPkFilename = "transition.proving_key"
const transitionCsFilename = "transition.constraint_system"

type ProofService struct {
	instantiationCs constraint.ConstraintSystem
	instantiationPk groth16.ProvingKey
	transitionCs    constraint.ConstraintSystem
	transitionPk    groth16.ProvingKey
}

var instantiationCircuit circuit.InstantiationCircuit
var transitionCircuit circuit.TransitionCircuit

func NewProofService() ProofService {
	instantiationCsPath := getFolderPath() + "/" + instantiationCsFilename
	instantiationPkPath := getFolderPath() + "/" + instantiationPkFilename
	transitionCsPath := getFolderPath() + "/" + transitionCsFilename
	transitionPkPath := getFolderPath() + "/" + transitionPkFilename
	instantiationCs := importConstraintSystem(&instantiationCircuit, instantiationCsPath)
	transitionCs := importConstraintSystem(&transitionCircuit, transitionCsPath)
	instantiationPk := importProvingKey(instantiationPkPath, instantiationCs)
	transitionPk := importProvingKey(transitionPkPath, transitionCs)
	return ProofService{
		instantiationCs: instantiationCs,
		instantiationPk: instantiationPk,
		transitionCs:    transitionCs,
		transitionPk:    transitionPk,
	}
}

func importConstraintSystem(circuit frontend.Circuit, path string) constraint.ConstraintSystem {
	cs, err := readConstraintSystem(path)
	if err != nil {
		cs = compileCircuit(circuit, path)
	}
	return cs
}

func readConstraintSystem(path string) (constraint.ConstraintSystem, error) {
	var cs constraint.ConstraintSystem
	err := readFile(cs, path)
	return cs, err
}

func compileCircuit(circuit frontend.Circuit, path string) constraint.ConstraintSystem {
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	check(err)
	writeFile(cs, path)
	return cs
}

func importProvingKey(path string, cs constraint.ConstraintSystem) groth16.ProvingKey {
	pk, err := readProvingKey(path)
	if err != nil {
		pk = generateProvingKey(cs, path)
	}
	return pk
}

func readProvingKey(path string) (groth16.ProvingKey, error) {
	var pk groth16.ProvingKey
	err := readFile(pk, path)
	return pk, err
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

func getFolderPath() string {
	return "./files"
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
