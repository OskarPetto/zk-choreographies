package parameters

import (
	"bytes"
	"crypto/rand"
	"execution-service/domain"
	"execution-service/files"
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
	err := files.ReadPrivateFile(byteBuffer, filename)
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
	files.WritePrivateFile(byteBuffer, filename)
	return privateKey
}

func (service *SignatureParameters) GetPrivateKeyForIdentity(identityId domain.IdentityId) (*eddsa.PrivateKey, error) {
	if int(identityId) >= len(service.privateKeys) {
		return nil, fmt.Errorf("private key does not exist for identity %d", identityId)
	}
	return service.privateKeys[identityId], nil
}

func (service *SignatureParameters) GetPublicKeys(count int) []domain.PublicKey {
	publicKeys := make([]domain.PublicKey, count)
	for i := 0; i < count; i++ {
		publicKey, err := service.GetPrivateKeyForIdentity(uint(i))
		if err != nil {
			break
		}
		publicKeys[i] = domain.PublicKey{
			Value: publicKey.PublicKey.Bytes(),
		}
	}
	return publicKeys
}
