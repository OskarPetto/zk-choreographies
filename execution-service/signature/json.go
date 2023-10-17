package signature

import (
	"execution-service/domain"
	"execution-service/utils"
)

type SignatureJson struct {
	Value     string `json:"value"`
	Instance  string `json:"instance"`
	PublicKey string `json:"publicKey"`
}

func ToJson(signature domain.Signature) SignatureJson {
	return SignatureJson{
		Value:     utils.BytesToString(signature.Value),
		Instance:  utils.BytesToString(signature.Instance.Value[:]),
		PublicKey: utils.BytesToString(signature.PublicKey.Value),
	}
}

func (signatureJson *SignatureJson) ToSignature() (domain.Signature, error) {
	value, err := utils.StringToBytes(signatureJson.Value)
	if err != nil {
		return domain.Signature{}, err
	}
	instance, err := utils.StringToBytes(signatureJson.Instance)
	if err != nil {
		return domain.Signature{}, err
	}
	publicKey, err := utils.StringToBytes(signatureJson.PublicKey)
	if err != nil {
		return domain.Signature{}, err
	}
	return domain.Signature{
		Value: value,
		Instance: domain.Hash{
			Value: [domain.HashSize]byte(instance),
		},
		PublicKey: domain.PublicKey{
			Value: publicKey,
		},
	}, nil
}
