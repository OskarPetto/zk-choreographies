package instance

import (
	"fmt"
	"proof-service/petri_net"
	"proof-service/workflow"
	"slices"
)

var AllowedTokenCounts = []int{0, 1}

type Instance struct {
	PlaceCount  uint8
	TokenCounts [petri_net.MaxPlaceCount]int8
}

func FromWorkflowInstance(instance workflow.Instance) (Instance, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > petri_net.MaxPlaceCount {
		return Instance{}, fmt.Errorf("instance '%s' is too large", instance.Id)
	}
	var tokenCounts [petri_net.MaxPlaceCount]int8
	for i := 0; i < placeCount; i++ {
		tokenCount := instance.TokenCounts[i]
		if !slices.Contains(AllowedTokenCounts, tokenCount) {
			return Instance{}, fmt.Errorf("tokenCount of instance '%s' at index %d is not allowed", instance.Id, i)
		}
		tokenCounts[i] = int8(instance.TokenCounts[i])
	}
	return Instance{
		PlaceCount:  uint8(placeCount),
		TokenCounts: tokenCounts,
	}, nil
}

func SerializeInstance(instance Instance) []byte {
	var bytes = make([]byte, petri_net.MaxPlaceCount+1)
	bytes[0] = instance.PlaceCount
	for i := 0; i < int(instance.PlaceCount); i++ {
		bytes[i+1] = byte(instance.TokenCounts[i])
	}
	return bytes
}
