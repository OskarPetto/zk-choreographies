package domain

import (
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
	MessageHashes [MaxMessageCount][HashSize]byte
	CreatedAt     int64
}

func (instance *Instance) Id() InstanceId {
	return instance.Hash.String()
}

func (instance Instance) SetMessageHash(messageId MessageId, messageHash Hash) Instance {
	instance.MessageHashes[messageId] = messageHash.Value
	instance.CreatedAt = time.Now().Unix()
	instance.ComputeHash()
	return instance
}

func (instance Instance) UpdateTokenCounts(transition Transition, input ConstraintInput) (Instance, error) {
	if !isTransitionExecutable(instance, transition, input) {
		return Instance{}, fmt.Errorf("transition %s is not executable", transition.Id)
	}

	instance.updateTokenCounts(transition)
	instance.CreatedAt = time.Now().Unix()
	instance.ComputeHash()
	return instance, nil
}

func (instance *Instance) updateTokenCounts(transition Transition) {
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
}

func isTransitionExecutable(instance Instance, transition Transition, input ConstraintInput) bool {
	if !transition.IsTransition {
		return false
	}
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
	return instance.EvaluateConstraint(transition.Constraint, input)
}
