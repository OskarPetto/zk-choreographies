package hash

import (
	"execution-service/domain"
	"execution-service/utils"
)

type SaltedHashJson struct {
	Hash string `json:"hash"`
	Salt string `json:"salt"`
}

func ToJson(hash domain.SaltedHash) SaltedHashJson {
	return SaltedHashJson{
		Hash: utils.BytesToString(hash.Hash.Value[:]),
		Salt: utils.BytesToString(hash.Salt[:]),
	}
}

func (hash *SaltedHashJson) ToHash() (domain.SaltedHash, error) {
	if hash.Hash == "" && hash.Salt == "" {
		return domain.SaltedHash{}, nil
	}
	value, err := utils.StringToBytes(hash.Hash)
	if err != nil {
		return domain.SaltedHash{}, err
	}
	salt, err := utils.StringToBytes(hash.Salt)
	if err != nil {
		return domain.SaltedHash{}, err
	}
	valueFixed := [domain.HashSize]byte(value)
	saltFixed := [domain.SaltSize]byte(salt)
	return domain.SaltedHash{
		Hash: domain.Hash{
			Value: valueFixed,
		},
		Salt: saltFixed,
	}, nil
}
