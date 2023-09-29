package parameters

import (
	"bytes"
	"crypto/rand"
	"execution-service/domain"
	"execution-service/utils"
	"fmt"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

type SignatureParameters struct {
	privateKeys []*eddsa.PrivateKey
}

func NewSignatureParameters() SignatureParameters {
	signaturePrivateKeys := make([]*eddsa.PrivateKey, domain.IdentityCount)
	for i := 0; i < domain.IdentityCount; i++ {
		signaturePrivateKeys[i] = importSignaturePrivateKey(fmt.Sprintf("identity%d.private_key", i))
	}
	return SignatureParameters{
		signaturePrivateKeys,
	}
}

func importSignaturePrivateKey(filename string) *eddsa.PrivateKey {
	var pk eddsa.PrivateKey
	byteBuffer := new(bytes.Buffer)
	err := readPrivateFile(byteBuffer, filename)
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
	writePrivateFile(byteBuffer, filename)
	return privateKey
}

func (service *SignatureParameters) GetPrivateKeyForIdentity(identityId domain.IdentityId) *eddsa.PrivateKey {
	return service.privateKeys[identityId]
}

func (service *SignatureParameters) GetPublicKeys(count int) []domain.PublicKey {
	publicKeys := make([]domain.PublicKey, count)
	for i := 0; i < count; i++ {
		publicKeys[i] = domain.PublicKey{
			Value: service.GetPrivateKeyForIdentity(uint(i)).PublicKey.Bytes(),
		}
	}
	return publicKeys
}
