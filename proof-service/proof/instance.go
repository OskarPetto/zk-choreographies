package proof

import "proof-service/commitment"

const MaxPlaceCount = 100

type Instance struct {
	TokenCountsLength uint8
	TokenCounts       [MaxPlaceCount]int8
	Commitment        commitment.Commitment
}
