package domain

import (
	"fmt"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

const PublicKeySize = 32

type InstanceId = string

type Instance struct {
	Hash          SaltedHash
	Model         Hash
	TokenCounts   []int8
	PublicKeys    []PublicKey
	MessageHashes []Hash
	CreatedAt     int64
}

type PublicKey struct {
	Value []byte
}

const OutOfBoundsTokenCount = -100

func OutOfBoundsPublicKey() PublicKey {
	return PublicKey{
		Value: make([]byte, PublicKeySize),
	}
}

func NewPublicKey(eddsaPub eddsa.PublicKey) PublicKey {
	return PublicKey{
		Value: eddsaPub.Bytes(),
	}
}

func (instance *Instance) Id() InstanceId {
	return instance.Hash.String()
}

func (instance Instance) ExecuteTransition(transition Transition, input ConstraintInput, message *Message) (Instance, error) {
	if !isTransitionExecutable(instance, transition, input) {
		return Instance{}, fmt.Errorf("transition %s is not executable", transition.Id)
	}
	if transition.Message != EmptyMessageId && message == nil {
		return Instance{}, fmt.Errorf("transition %s requires a message", transition.Id)
	}
	if transition.Message == EmptyMessageId && message != nil {
		return Instance{}, fmt.Errorf("transition %s does not send any message", transition.Id)
	}
	instance.updateTokenCounts(transition)
	if message != nil {
		messageHash := message.Hash.Hash
		messageHashes := make([]Hash, len(instance.MessageHashes))
		copy(messageHashes[:], instance.MessageHashes[:])
		if transition.Message != EmptyMessageId {
			messageHashes[transition.Message] = messageHash
		}
		instance.MessageHashes = messageHashes
	}
	instance.CreatedAt = time.Now().Unix()
	instance.UpdateHash()
	return instance, nil
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

func (instance *Instance) FindMessageHashById(id ModelMessageId) Hash {
	return instance.MessageHashes[id]
}

func (instance *Instance) FindPublicKeyByParticipant(id ParticipantId) PublicKey {
	return instance.PublicKeys[id]
}
