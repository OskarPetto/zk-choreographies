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

var instantiationCircuit circuit.InstantiationCircuit

func TestInstantiation(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance1(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstantiation_InvalidSaltedHash(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance1(publicKey)
	instanceSaltedHash := crypto.HashInstance(testdata.GetPetriNet1Instance2(publicKey))
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidTokenCounts(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance2(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidSignature(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := signatureService.GetPublicKey()
	instance := testdata.GetPetriNet1Instance1(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(crypto.HashInstance(instance))
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}

func TestInstantiation_InvalidAuthorization(t *testing.T) {
	signatureService := crypto.NewSignatureService()
	publicKey := testdata.GetPublicKeys(2)[1]
	instance := testdata.GetPetriNet1Instance1(publicKey)
	instanceSaltedHash := crypto.HashInstance(instance)
	signature := signatureService.Sign(instanceSaltedHash)
	circuitInstance, _ := circuit.FromInstance(instance)

	petriNet, _ := circuit.FromPetriNet(testdata.GetPetriNet1())
	witness := circuit.InstantiationCircuit{
		Instance:   circuitInstance,
		SaltedHash: circuit.FromSaltedHash(instanceSaltedHash),
		Signature:  circuit.FromSignature(signature),
		PetriNet:   petriNet,
	}

	err := test.IsSolved(&instantiationCircuit, &witness, ecc.BN254.ScalarField())
	assert.NotNil(t, err)
}
