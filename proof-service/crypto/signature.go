package crypto

import (
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
)

type Signature struct {
	Value     []byte
	PublicKey eddsa.PublicKey
}

func Sign(commitment Commitment) Signature {
	privateKey := parameters.signaturePrivateKey
	signature, err := privateKey.Sign(commitment.Value, hash.MIMC_BN254.New())
	utils.PanicOnError(err)

	return Signature{
		Value:     signature,
		PublicKey: privateKey.PublicKey,
	}
}
