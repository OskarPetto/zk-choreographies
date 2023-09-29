package domain

import (
	"execution-service/utils"
	"fmt"
	"time"
)

const PublicKeySize = 32

type PublicKey struct {
	Value []byte
}

const OutOfBoundsTokenCount = -100

func OutOfBoundsPublicKey() PublicKey {
	return PublicKey{
		Value: make([]byte, PublicKeySize),
	}
}

type InstanceId = string

type Instance struct {
	Hash          Hash
	Model         ModelId
	TokenCounts   [MaxPlaceCount]int8
	PublicKeys    [MaxParticipantCount]PublicKey
	MessageHashes [MaxMessageCount]Hash
	CreatedAt     int64
}

func (instance *Instance) Id() InstanceId {
	return utils.BytesToString(instance.Hash.Value[:])
}

func (instance Instance) ExecuteTransition(transition Transition, message []byte) (Instance, error) {
	instance.updateMessageHash(transition.Message, message)
	err := instance.updateTokenCounts(transition)
	if err != nil {
		return Instance{}, err
	}
	instance.CreatedAt = time.Now().Unix()
	instance.ComputeHash()
	return instance, nil
}

func (instance *Instance) updateMessageHash(messageId MessageId, message []byte) {
	if messageId != EmptyMessageId && len(message) > 0 {
		messageHash := HashMessage(message)
		instance.MessageHashes[messageId] = messageHash
	}
}

func (instance *Instance) updateTokenCounts(transition Transition) error {
	if !isTransitionExecutable(instance, transition) {
		return fmt.Errorf("transition %s is not executable", transition.Id)
	}
	var tokenCounts [MaxPlaceCount]int8
	copy(tokenCounts[:], instance.TokenCounts[:])
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if incomingPlaceId != OutOfBoundsPlaceId {
			tokenCounts[incomingPlaceId] -= 1
		}
	}
	for _, outgoingPlaceId := range transition.OutgoingPlaces {
		if outgoingPlaceId != OutOfBoundsPlaceId {
			tokenCounts[outgoingPlaceId] += 1
		}
	}
	instance.TokenCounts = tokenCounts
	return nil
}

func isTransitionExecutable(instance *Instance, transition Transition) bool {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if incomingPlaceId == OutOfBoundsPlaceId {
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
