package commitment

const randomnessSize = 32

type CommitmentId = string

type Commitment struct {
	Id         CommitmentId
	Value      []byte
	Randomness []byte
}
