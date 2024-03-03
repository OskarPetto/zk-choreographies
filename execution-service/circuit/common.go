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

func checkModelHash(api frontend.API, model Model, instance Instance) error {
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
	mimc.Write(model.Salt)
	result := mimc.Sum()
	api.AssertIsEqual(result, instance.Model)
	return nil
}

func checkInstanceHash(api frontend.API, instance Instance) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(instance.Model)
	mimc.Write(instance.TokenCounts[:]...)
	mimc.Write(instance.PublicKeyRoot)
	mimc.Write(instance.MessageHashes[:]...)
	mimc.Write(instance.Hash.Salt)
	hash := mimc.Sum()
	api.AssertIsEqual(hash, instance.Hash.Hash)
	return nil
}

func checkAuthentication(api frontend.API, authentication Authentication, instance Instance) error {
	authentication.MerkleProof.CheckRootHash(api, instance.PublicKeyRoot)
	err := authentication.MerkleProof.VerifyProof(api)
	if err != nil {
		return err
	}
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

	return eddsa.Verify(curve, authentication.Signature, instance.Hash.Hash, authentication.PublicKey, &mimc)
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

func (merkleProof *MerkleProof) CheckRootHash(api frontend.API, hash frontend.Variable) {
	api.AssertIsEqual(merkleProof.MerkleProof.RootHash, hash)
}

func (merkleProof *MerkleProof) VerifyProof(api frontend.API) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	merkleProof.MerkleProof.VerifyProof(api, &mimc, merkleProof.Index)
	return nil
}
