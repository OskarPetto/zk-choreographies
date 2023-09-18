package proof

import (
	"proof-service/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type verifyingKey groth16.VerifyingKey

type ProofService struct {
	instantiationCs constraint.ConstraintSystem
	instantiationPk groth16.ProvingKey
	executionCs     constraint.ConstraintSystem
	executionPk     groth16.ProvingKey
}

var instantiationCircuit circuit.InstantiationCircuit
var executionCircuit circuit.ExecutionCircuit

func NewProofService() ProofService {
	instantiationCs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &instantiationCircuit)
	if err != nil {
		panic(err)
	}
	executionCs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &executionCircuit)
	if err != nil {
		panic(err)
	}

	return ProofService{
		instantiationCs: instantiationCs,
		instantiationPk: nil,
		executionCs:     executionCs,
		executionPk:     nil,
	}
}

func (service *ProofService) importKeys() error {
	return nil
}

// func (service *ProofService) ExportKeys() {
// 	pk, _, err := groth16.Setup(service.circuitService.InstantiationCircuitCs)
// 	if err != nil {
// 		panic(err)
// 	}
// 	file, err := os.Create("./proving_key")
// 	if err != nil {
// 		panic(err)
// 	}
// 	writer := bufio.NewWriter(file)
// 	bytesWritten, err := pk.WriteTo(writer)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("Bytes Written: %d\n", bytesWritten)
// 	writer.Flush()
// }
