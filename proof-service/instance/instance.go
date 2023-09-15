package instance

import (
	"fmt"
	"proof-service/petri_net"
	"slices"
)

var AllowedTokenCounts = []int{0, 1}

type Instance struct {
	Id          string
	TokenCounts []int
}

func ValidateInstance(instance Instance) error {
	placeCount := len(instance.TokenCounts)
	if placeCount > petri_net.MaxPlaceCount {
		return fmt.Errorf("instance '%s' is too large", instance.Id)
	}
	for i := 0; i < placeCount; i++ {
		tokenCount := instance.TokenCounts[i]
		if !slices.Contains(AllowedTokenCounts, tokenCount) {
			return fmt.Errorf("tokenCount of instance '%s' at index %d is not allowed", instance.Id, i)
		}
	}
	return nil
}

func SerializeInstance(instance Instance) []byte {
	var bytes = make([]byte, len(instance.TokenCounts))
	for i := 0; i < len(instance.TokenCounts); i++ {
		bytes[i] = byte(instance.TokenCounts[i])
	}
	return bytes
}
