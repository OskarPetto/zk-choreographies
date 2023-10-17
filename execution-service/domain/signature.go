package domain

import (
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
)

type IdentityId = uint

const IdentityCount = 3

type Signature struct {
	Value     []byte
	Instance  Hash
	PublicKey PublicKey
}

func (instance Instance) Sign(privateKey *eddsa.PrivateKey) Signature {
	signature, err := privateKey.Sign(instance.Hash.Hash.Value[:], hash.MIMC_BN254.New())
	utils.PanicOnError(err)
	publicKey := PublicKey{
		Value: privateKey.PublicKey.Bytes(),
	}
	return Signature{
		Value:     signature,
		Instance:  instance.Hash.Hash,
		PublicKey: publicKey,
	}
}

func (signature Signature) Verify() bool {
	var publicKey eddsa.PublicKey
	publicKey.SetBytes(signature.PublicKey.Value)
	doesVerify, err := publicKey.Verify(signature.Value, signature.Instance.Value[:], hash.MIMC_BN254.New())
	utils.PanicOnError(err)
	return doesVerify
}
