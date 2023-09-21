package circuit

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

func checkInstanceSaltedHash(api frontend.API, instance Instance, saltedHash SaltedHash) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(instance.TokenCounts[:]...)
	for _, publicKey := range instance.PublicKeys {
		mimc.Write(publicKey.A.X)
		mimc.Write(publicKey.A.Y)
	}
	mimc.Write(saltedHash.Salt)
	hash := mimc.Sum()
	api.AssertIsEqual(hash, saltedHash.Value)
	return nil
}

func checkSignature(api frontend.API, signature Signature, saltedHash SaltedHash) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	return eddsa.Verify(curve, signature.Value, saltedHash.Value, signature.PublicKey, &mimc)
}

func findParticipantId(api frontend.API, signature Signature, instance Instance) frontend.Variable {
	var participantId frontend.Variable = -1
	for i, publicKey := range instance.PublicKeys {
		participantId = api.Select(publicKeyEquals(api, publicKey, signature.PublicKey), i, participantId)
	}
	api.AssertIsDifferent(participantId, -1)
	return participantId
}

func publicKeyEquals(api frontend.API, a, b eddsa.PublicKey) frontend.Variable {
	return api.And(equals(api, a.A.X, b.A.X), equals(api, a.A.Y, b.A.Y))
}

func equals(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.IsZero(api.Sub(a, b))
}
