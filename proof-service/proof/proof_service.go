package proof

import (
	"bufio"
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
const executionPkFilename = "execution.proving_key"
const executionCsFilename = "execution.constraint_system"

type ProofService struct {
	instantiationCs *constraint.ConstraintSystem
	instantiationPk *groth16.ProvingKey
	executionCs     *constraint.ConstraintSystem
	executionPk     *groth16.ProvingKey
}

var instantiationCircuit circuit.InstantiationCircuit
var executionCircuit circuit.ExecutionCircuit

func NewProofService() ProofService {
	instantiationCsPath := getFolderPath() + "/" + instantiationCsFilename
	instantiationCs := importConstraintSystem(&instantiationCircuit, instantiationCsPath)
	instantiationPkPath := getFolderPath() + "/" + instantiationPkFilename
	instantiationPk := importProvingKey(instantiationPkPath, instantiationCs)
	executionCsPath := getFolderPath() + "/" + executionCsFilename
	executionCs := importConstraintSystem(&executionCircuit, executionCsPath)
	executionPkPath := getFolderPath() + "/" + executionPkFilename
	executionPk := importProvingKey(executionPkPath, executionCs)
	return ProofService{
		instantiationCs: instantiationCs,
		instantiationPk: instantiationPk,
		executionCs:     executionCs,
		executionPk:     executionPk,
	}
}

func importConstraintSystem(circuit frontend.Circuit, path string) *constraint.ConstraintSystem {
	cs, err := readConstraintSystem(path)
	if err != nil {
		cs = compileCircuit(circuit, path)
	}
	return cs
}

func importProvingKey(path string, cs *constraint.ConstraintSystem) *groth16.ProvingKey {
	pk, err := readProvingKey(path)
	if err != nil {
		pk = generateProvingKey(cs, path)
	}
	return pk
}

func readConstraintSystem(path string) (*constraint.ConstraintSystem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var cs constraint.ConstraintSystem
	cs.ReadFrom(file)
	return &cs, nil
}

func readProvingKey(path string) (*groth16.ProvingKey, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var pk groth16.ProvingKey
	pk.ReadFrom(file)
	return &pk, nil
}

func compileCircuit(circuit frontend.Circuit, path string) *constraint.ConstraintSystem {
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, circuit)
	check(err)
	writeFile(cs, path)
	return &cs
}

func generateProvingKey(cs *constraint.ConstraintSystem, path string) *groth16.ProvingKey {
	pk, _, err := groth16.Setup(*cs)
	check(err)
	writeFile(pk, path)
	return &pk
}

func writeFile(writeable io.WriterTo, path string) {
	file, err := os.Create(path)
	check(err)
	defer file.Close()
	writer := bufio.NewWriter(file)
	bytesWritten, err := writeable.WriteTo(writer)
	check(err)
	fmt.Printf("Generated file of size %d in %s\n", bytesWritten, path)
	writer.Flush()
}

func getFolderPath() string {
	return "/home/opetto/uni/zk-choreographies/files"
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
