package domain

import (
	"bytes"
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
	"github.com/consensys/gnark-crypto/hash"
)

type IdentityId = uint

const IdentityCount = MaxParticipantCount

type Signature struct {
	Value       []byte
	PublicKey   PublicKey
	Participant ParticipantId
}

func (instance Instance) Sign(privateKey *eddsa.PrivateKey) Signature {
	signature, err := privateKey.Sign(instance.Hash.Value[:], hash.MIMC_BN254.New())
	utils.PanicOnError(err)
	publicKey := PublicKey{
		Value: privateKey.PublicKey.Bytes(),
	}
	var participantId ParticipantId = 0
	for i, instancePublicKey := range instance.PublicKeys {
		if bytes.Equal(publicKey.Value, instancePublicKey.Value) {
			participantId = ParticipantId(i)
			break
		}
	}
	return Signature{
		Value:       signature,
		PublicKey:   publicKey,
		Participant: participantId,
	}
}
