package commitment

const RandomnessSize = 32

type CommitmentId = string

type Commitment struct {
	Id         CommitmentId
	Value      [32]byte
	Randomness []byte
}
