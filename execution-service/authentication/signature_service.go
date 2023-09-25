package authentication

import (
	"proof-service/authentication/parameters"
	"proof-service/domain"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/hash"
)

type Signature struct {
	Value     []byte
	PublicKey domain.PublicKey
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

func (service *SignatureService) Sign(instance domain.Instance) Signature {
	privateKey := service.signatureParameters.SignaturePrivateKey
	signature, err := privateKey.Sign(instance.Hash, hash.MIMC_BN254.New())
	utils.PanicOnError(err)

	return Signature{
		Value:     signature,
		PublicKey: service.GetPublicKey(),
	}
}

func (service *SignatureService) GetPublicKey() domain.PublicKey {
	privateKey := service.signatureParameters.SignaturePrivateKey
	return domain.PublicKey{
		Value: privateKey.PublicKey.Bytes(),
	}
}
