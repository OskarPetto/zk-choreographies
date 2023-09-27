package authentication

import (
	"execution-service/authentication/parameters"
	"execution-service/domain"
	"execution-service/utils"

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
	signature, err := privateKey.Sign(instance.Hash.Value[:], hash.MIMC_BN254.New())
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
