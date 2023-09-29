package domain

import (
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
)

type IdentityId = uint

const IdentityCount = MaxParticipantCount

type Signature struct {
	Value     []byte
	PublicKey PublicKey
}

func (instance Instance) Sign(privateKey *eddsa.PrivateKey) Signature {
	signature, err := privateKey.Sign(instance.Hash.Value[:], hash.MIMC_BN254.New())
	utils.PanicOnError(err)
	return Signature{
		Value: signature,
		PublicKey: PublicKey{
			Value: privateKey.PublicKey.Bytes(),
		},
	}
}
