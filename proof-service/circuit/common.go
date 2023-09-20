package circuit

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

func checkCommitment(api frontend.API, instance Instance, commitment Commitment) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(instance.PlaceCount)
	mimc.Write(instance.TokenCounts[:]...)
	mimc.Write(commitment.Randomness)
	hash := mimc.Sum()
	api.AssertIsEqual(hash, commitment.Value)
	return nil
}

func equals(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.IsZero(api.Sub(a, b))
}

func checkSignature(api frontend.API, signature Signature, commitment Commitment) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	return eddsa.Verify(curve, signature.Value, commitment.Value, signature.PublicKey, &mimc)
}
