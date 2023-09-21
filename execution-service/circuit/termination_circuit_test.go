package circuit_test

import (
	"proof-service/circuit"
	"proof-service/crypto"
	"proof-service/testdata"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/test"
	"github.com/stretchr/testify/assert"
)

var terminationCircuit circuit.TerminationCircuit

func TestTermination(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance3(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TerminationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestTermination_InvalidSaltedHash(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance3(publicKey)
	instanceSaltedHash := crypto.HashInstance(testdata.GetPetriNet1Instance1(publicKey))
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TerminationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidTokenCounts(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance2(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TerminationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&terminationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidSignature(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance3(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(crypto.HashInstance(instance))
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TerminationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestTermination_InvalidAuthorization(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := testdata.GetPublicKey1()
	instance := testdata.GetPetriNet1Instance3(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.TerminationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
