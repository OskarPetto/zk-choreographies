package domain

import (
	"crypto/sha256"
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

func (instance Instance) ExecuteTransition(transition Transition, message []byte) (Instance, error) {
	if len(message) > 0 {
		err := instance.storeMessageHash(transition.Message, message)
		if err != nil {
			return Instance{}, err
		}
	}
	err := instance.executeTransition(transition)
	if err != nil {
		return Instance{}, err
	}
	return instance, nil
}

func (instance *Instance) storeMessageHash(messageId MessageId, message []byte) error {
	if messageId >= MaxMessageCount {
		return fmt.Errorf("messageId %d is not valid", messageId)
	}
	instance.MessageHashes[messageId] = MessageHash{
		Value: sha256.Sum256(message),
	}
	return nil
}

func (instance *Instance) executeTransition(transition Transition) error {
	if !isTransitionExecutable(instance, transition) {
		return fmt.Errorf("transition %s is not executable", transition.Id)
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
	return nil
}

func isTransitionExecutable(instance *Instance, transition Transition) bool {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if instance.TokenCounts[incomingPlaceId] < 1 {
			return false
		}
	}
	return true
}
