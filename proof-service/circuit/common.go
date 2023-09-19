package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

func checkCommitment(api frontend.API, instance Instance, comm Commitment) error {
	hasher, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	hasher.Write(instance.PlaceCount)
	hasher.Write(instance.TokenCounts[:]...)
	hasher.Write(comm.Randomness)
	hash := hasher.Sum()
	api.AssertIsEqual(hash, comm.Value)
	return nil
}

func equals(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.IsZero(api.Sub(a, b))
}
