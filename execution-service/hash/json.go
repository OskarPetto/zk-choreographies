package hash

import (
	"execution-service/domain"
	"execution-service/utils"
)

type HashJson struct {
	Value string `json:"value"`
	Salt  string `json:"salt"`
}

func ToJson(hash domain.Hash) HashJson {
	return HashJson{
		Value: utils.BytesToString(hash.Value[:]),
		Salt:  utils.BytesToString(hash.Salt[:]),
	}
}

func (hash *HashJson) ToHash() (domain.Hash, error) {
	if hash.Value == "" && hash.Salt == "" {
		return domain.EmptyHash(), nil
	}
	value, err := utils.StringToBytes(hash.Value)
	if err != nil {
		return domain.Hash{}, err
	}
	salt, err := utils.StringToBytes(hash.Salt)
	if err != nil {
		return domain.Hash{}, err
	}
	valueFixed := [domain.HashSize]byte(value)
	saltFixed := [domain.SaltSize]byte(salt)
	return domain.Hash{
		Value: valueFixed,
		Salt:  saltFixed,
	}, nil
}
