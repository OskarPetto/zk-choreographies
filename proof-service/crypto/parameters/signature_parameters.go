package parameters

import (
	"bytes"
	"crypto/rand"
	"proof-service/file"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

const signaturePrivateKeyFilename = "signature.private_key"

type SignatureParameters struct {
	isLoaded            bool
	SignaturePrivateKey *eddsa.PrivateKey
}

var signatureParameters SignatureParameters

func NewSignatureParameters() SignatureParameters {
	if !signatureParameters.isLoaded {
		signaturePrivateKey := importSignaturePrivateKey(signaturePrivateKeyFilename)
		signatureParameters = SignatureParameters{
			true,
			signaturePrivateKey,
		}
	}
	return signatureParameters
}

func importSignaturePrivateKey(filename string) *eddsa.PrivateKey {
	var pk eddsa.PrivateKey
	byteBuffer := new(bytes.Buffer)
	err := file.ReadFile(byteBuffer, filename)
	if err != nil {
		pk = *generateSignaturePrivateKey(filename)
	} else {
		pk.SetBytes(byteBuffer.Bytes())
	}
	return &pk
}

func generateSignaturePrivateKey(filename string) *eddsa.PrivateKey {
	privateKey, err := eddsa.GenerateKey(rand.Reader)
	utils.PanicOnError(err)
	byteBuffer := bytes.NewBuffer(privateKey.Bytes())
	file.WriteFile(byteBuffer, filename)
	return privateKey
}
