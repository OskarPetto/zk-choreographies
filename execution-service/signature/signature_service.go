//go:build wireinject
// +build wireinject

package signature

import (
	"execution-service/domain"
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/hash"
	"github.com/google/wire"
)

type SignatureService struct {
	signatureParameters SignatureParameters
}

func InitializeSignatureService() SignatureService {
	wire.Build(NewSignatureService, NewSignatureParameters)
	return SignatureService{}
}

func NewSignatureService(signatureParameters SignatureParameters) SignatureService {
	return SignatureService{
		signatureParameters: signatureParameters,
	}
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
