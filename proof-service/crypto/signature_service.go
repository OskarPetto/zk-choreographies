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
	signatureParameters parameters.SignatureParameters
}

func NewSignatureService() SignatureService {
	return SignatureService{
		signatureParameters: parameters.NewSignatureParameters(),
	}
}

func (service *SignatureService) Sign(commitment Commitment) Signature {
	privateKey := service.signatureParameters.SignaturePrivateKey
	signature, err := privateKey.Sign(commitment.Value, hash.MIMC_BN254.New())
	utils.PanicOnError(err)

	return Signature{
		Value:     signature,
		PublicKey: privateKey.PublicKey.Bytes(),
	}
}
