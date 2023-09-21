package domain

type Instance struct {
	Id          string
	TokenCounts []int
	PublicKeys  [][]byte
}
