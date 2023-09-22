package crypto

import (
	"proof-service/crypto/parameters"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/hash"
)

type Signature struct {
	Value     []byte
	PublicKey []byte
}

type SignatureService struct {
	isLoaded            bool
	signatureParameters parameters.SignatureParameters
}

var signatureService SignatureService

func NewSignatureService() SignatureService {
	if !signatureService.isLoaded {
		signatureService = SignatureService{
			isLoaded:            true,
			signatureParameters: parameters.LoadSignatureParameters(),
		}
	}
	return signatureService
}

func (service *SignatureService) Sign(saltedHash SaltedHash) Signature {
	privateKey := service.signatureParameters.SignaturePrivateKey
	signature, err := privateKey.Sign(saltedHash.Value, hash.MIMC_BN254.New())
	utils.PanicOnError(err)

	return Signature{
		Value:     signature,
		PublicKey: service.GetPublicKey(),
	}
}

func (service *SignatureService) GetPublicKey() []byte {
	privateKey := service.signatureParameters.SignaturePrivateKey
	return privateKey.PublicKey.Bytes()
}
