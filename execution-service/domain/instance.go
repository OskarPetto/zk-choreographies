package domain

import (
	"fmt"
)

const MessageHashLength = 32

type MessageHash struct {
	Value [MessageHashLength]byte
}

type PublicKey struct {
	Value []byte
}

type Instance struct {
	Hash          []byte
	TokenCounts   []int
	MessageHashes []MessageHash
	PublicKeys    []PublicKey
	Salt          []byte
}

func (instance Instance) ExecuteTransition(transition Transition) (Instance, error) {
	if !isTransitionExecutable(instance, transition) {
		return instance, fmt.Errorf("transition %s is not executable", transition.Id)
	}
	tokenCounts := make([]int, len(instance.TokenCounts))
	copy(tokenCounts, instance.TokenCounts)
	for _, incomingPlaceId := range transition.IncomingPlaces {
		tokenCounts[incomingPlaceId] -= 1
	}
	for _, outgoingPlaceId := range transition.OutgoingPlaces {
		tokenCounts[outgoingPlaceId] += 1
	}
	instance.TokenCounts = tokenCounts
	instance.ComputeHash()
	return instance, nil
}

func isTransitionExecutable(instance Instance, transition Transition) bool {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if instance.TokenCounts[incomingPlaceId] < 1 {
			return false
		}
	}
	return true
}
