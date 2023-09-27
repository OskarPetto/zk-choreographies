package domain

import (
	"fmt"
)

const PublicKeySize = 32

type PublicKey struct {
	Value []byte
}

const DefaultTokenCount = 0

var DefaultPublicKey = PublicKey{
	Value: make([]byte, PublicKeySize),
}

type InstanceId = string

type Instance struct {
	Id            InstanceId
	Model         ModelId
	Hash          Hash
	TokenCounts   [MaxPlaceCount]int8
	PublicKeys    [MaxParticipantCount]PublicKey
	MessageHashes [MaxMessageCount]Hash
}

func (instance Instance) ExecuteTransition(transition Transition) (Instance, error) {
	err := instance.executeTransition(transition)
	if err != nil {
		return Instance{}, err
	}
	return instance, nil
}

func (instance Instance) ExecuteTransitionWithMessage(transition Transition, message []byte) (Instance, error) {
	instance.storeMessageHash(transition.Message, message)
	err := instance.executeTransition(transition)
	if err != nil {
		return Instance{}, err
	}
	return instance, nil
}

func (instance *Instance) storeMessageHash(messageId MessageId, message []byte) {
	messageHash := HashMessage(message)
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
		if int(incomingPlaceId) > len(instance.TokenCounts) {
			return false
		}

		if instance.TokenCounts[incomingPlaceId] < 1 {
			return false
		}
	}
	return true
}
