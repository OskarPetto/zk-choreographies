package commitment

const RandomnessSize = 32
const CommitmentSize = 32

type CommitmentId = string

type Commitment struct {
	Id         CommitmentId
	Value      [CommitmentSize]byte
	Randomness [RandomnessSize]byte
}
