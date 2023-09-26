package domain

import (
	"fmt"
)

const MessageHashSize = 32
const PublicKeySize = 32

type MessageHash struct {
	Value [MessageHashSize]byte
}

type PublicKey struct {
	Value []byte
}

const DefaultTokenCount = 0

var DefaultPublicKey = PublicKey{
	Value: make([]byte, PublicKeySize),
}

var DefaultMessageHash = MessageHash{
	Value: [MessageHashSize]byte{},
}

type Instance struct {
	Hash          []byte
	TokenCounts   [MaxPlaceCount]int8
	PublicKeys    [MaxParticipantCount]PublicKey
	MessageHashes [MaxMessageCount]MessageHash
	Salt          []byte
}

func (instance Instance) ExecuteTransition(transition Transition, messageHash MessageHash) (Instance, error) {
	instance.storeMessageHash(transition.Message, messageHash)
	err := instance.executeTransition(transition)
	if err != nil {
		return Instance{}, err
	}
	return instance, nil
}

func (instance *Instance) storeMessageHash(messageId MessageId, messageHash MessageHash) {
	if messageId != DefaultMessageId {
		instance.MessageHashes[messageId] = messageHash
	}
}

func (instance *Instance) executeTransition(transition Transition) error {
	if !isTransitionExecutable(instance, transition) {
		return fmt.Errorf("transition %s is not executable", transition.Id)
	}
	var tokenCounts [MaxPlaceCount]int8
	copy(tokenCounts[:], instance.TokenCounts[:])
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if incomingPlaceId != DefaultPlaceId {
			tokenCounts[incomingPlaceId] -= 1
		}
	}
	for _, outgoingPlaceId := range transition.OutgoingPlaces {
		if outgoingPlaceId != DefaultPlaceId {
			tokenCounts[outgoingPlaceId] += 1
		}
	}
	instance.TokenCounts = tokenCounts
	instance.ComputeHash()
	return nil
}

func isTransitionExecutable(instance *Instance, transition Transition) bool {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if incomingPlaceId == DefaultPlaceId {
			break
		}
		if instance.TokenCounts[incomingPlaceId] < 1 {
			return false
		}
	}
	return true
}
