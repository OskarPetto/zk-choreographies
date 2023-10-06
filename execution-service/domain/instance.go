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
	TokenCounts   []int8
	PublicKeys    []PublicKey
	MessageHashes [][HashSize]byte
	CreatedAt     int64
}

func (instance *Instance) Id() InstanceId {
	return instance.Hash.String()
}

func (instance Instance) ExecuteTransition(transition Transition, input ConstraintInput, messageHash Hash) (Instance, error) {
	if !isTransitionExecutable(instance, transition, input) {
		return Instance{}, fmt.Errorf("transition %s is not executable", transition.Id)
	}
	instance.updateMessageHashes(transition, messageHash)
	instance.updateTokenCounts(transition)
	instance.CreatedAt = time.Now().Unix()
	instance.ComputeHash()
	return instance, nil
}

func (instance *Instance) updateMessageHashes(transition Transition, messageHash Hash) {
	messageHashes := make([][HashSize]byte, len(instance.MessageHashes))
	copy(messageHashes[:], instance.MessageHashes[:])
	if transition.Message != EmptyMessageId {
		messageHashes[transition.Message] = messageHash.Value
	}
	instance.MessageHashes = messageHashes
}

func (instance *Instance) updateTokenCounts(transition Transition) {
	tokenCounts := make([]int8, len(instance.TokenCounts))
	copy(tokenCounts[:], instance.TokenCounts[:])
	for _, incomingPlaceId := range transition.IncomingPlaces {
		tokenCounts[incomingPlaceId] -= 1
	}
	for _, outgoingPlaceId := range transition.OutgoingPlaces {
		tokenCounts[outgoingPlaceId] += 1
	}
	instance.TokenCounts = tokenCounts
}

func isTransitionExecutable(instance Instance, transition Transition, input ConstraintInput) bool {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if instance.TokenCounts[incomingPlaceId] < 1 {
			return false
		}
	}
	return instance.EvaluateConstraint(transition.Constraint, input)
}
