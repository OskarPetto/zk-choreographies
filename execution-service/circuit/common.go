package circuit

import (
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const emptyMessageHash = 0
const outOfBoundsMessageHash = 1

func checkModelHash(api frontend.API, model Model) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(model.PlaceCount)
	mimc.Write(model.ParticipantCount)
	mimc.Write(model.MessageCount)
	for _, startPlace := range model.StartPlaces {
		mimc.Write(startPlace)
	}
	mimc.Write(model.EndPlaceRoot)
	mimc.Write(model.TransitionRoot)
	mimc.Write(model.Hash.Salt)
	result := mimc.Sum()
	api.AssertIsEqual(result, model.Hash.Value)
	return nil
}

func checkInstanceHash(api frontend.API, instance Instance) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(instance.TokenCounts[:]...)
	mimc.Write(instance.PublicKeyRoot)
	mimc.Write(instance.MessageHashes[:]...)
	mimc.Write(instance.Hash.Salt)
	hash := mimc.Sum()
	api.AssertIsEqual(hash, instance.Hash.Value)
	return nil
}

func checkAuthentication(api frontend.API, authentication Authentication, instance Instance) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	authentication.MerkleProof.CheckRootHash(api, instance.PublicKeyRoot)
	authentication.MerkleProof.VerifyProof(api, mimc)
	checkPublicKeyHash(api, authentication.MerkleProof.MerkleProof.Path[0], authentication.PublicKey)
	return checkSignature(api, authentication, instance)
}

func checkSignature(api frontend.API, authentication Authentication, instance Instance) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	return eddsa.Verify(curve, authentication.Signature, instance.Hash.Value, authentication.PublicKey, &mimc)
}

func checkPublicKeyHash(api frontend.API, hash frontend.Variable, publicKey eddsa.PublicKey) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(publicKey.A.X)
	mimc.Write(publicKey.A.Y)
	api.AssertIsEqual(hash, mimc.Sum())
	return nil
}

func equals(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.IsZero(api.Sub(a, b))
}
