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
	SignaturePrivateKey *eddsa.PrivateKey
}

func LoadSignatureParameters() SignatureParameters {
	signaturePrivateKey := importSignaturePrivateKey(signaturePrivateKeyFilename)
	return SignatureParameters{
		signaturePrivateKey,
	}
}

func importSignaturePrivateKey(filename string) *eddsa.PrivateKey {
	var pk eddsa.PrivateKey
	byteBuffer := new(bytes.Buffer)
	err := file.ReadPrivateFile(byteBuffer, filename)
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
	file.WritePrivateFile(byteBuffer, filename)
	return privateKey
}
