package circuit

import (
	"proof-service/commitment"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/uints"
)

func checkCommitment(api frontend.API, uapi *uints.BinaryField[uints.U32], instance Instance, comm Commitment) error {
	hasher, err := sha2.New(api)
	if err != nil {
		return err
	}
	hasher.Write([]uints.U8{instance.PlaceCount})
	hasher.Write(instance.TokenCounts[:])
	hasher.Write(comm.Randomness[:])
	hash := hasher.Sum()
	for i := 0; i < commitment.CommitmentSize; i++ {
		uapi.ByteAssertEq(comm.Value[i], hash[i])
	}
	return nil
}

func equals(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.IsZero(api.Sub(a, b))
}
