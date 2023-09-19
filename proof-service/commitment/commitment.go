package commitment

const RandomnessSize = 32

type CommitmentId = string

type Commitment struct {
	Id         CommitmentId
	Value      []byte
	Randomness [RandomnessSize]byte
}
