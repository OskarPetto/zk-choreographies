package domain

import (
	"bytes"
	"fmt"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

const PublicKeySize = 32

type InstanceId = string

type Instance struct {
	SaltedHash    SaltedHash
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
	return instance.SaltedHash.String()
}

func (instance Instance) ExecuteTransition(transition Transition, input ConstraintInput, initiatingMessage *Message, respondingMessage *Message) (Instance, error) {
	err := validateTransitionExecutable(instance, transition, input)
	if err != nil {
		return Instance{}, err
	}
	if transition.InitiatingMessage != EmptyMessageId && initiatingMessage == nil {
		return Instance{}, fmt.Errorf("transition %s requires a initiating message", transition.Id)
	}
	if transition.InitiatingMessage == EmptyMessageId && initiatingMessage != nil {
		return Instance{}, fmt.Errorf("transition %s does not send any initiating message", transition.Id)
	}
	if transition.RespondingMessage != EmptyMessageId && respondingMessage == nil {
		return Instance{}, fmt.Errorf("transition %s requires a responding message", transition.Id)
	}
	if transition.RespondingMessage == EmptyMessageId && respondingMessage != nil {
		return Instance{}, fmt.Errorf("transition %s does not send any responding message", transition.Id)
	}
	instance.updateTokenCounts(transition)
	if initiatingMessage != nil {
		messageHashes := make([]Hash, len(instance.MessageHashes))
		copy(messageHashes[:], instance.MessageHashes[:])
		if transition.InitiatingMessage != EmptyMessageId {
			initiatingMessageHash := initiatingMessage.Hash.Hash
			messageHashes[transition.InitiatingMessage] = initiatingMessageHash
		}
		if transition.RespondingMessage != EmptyMessageId {
			respondingMessageHash := respondingMessage.Hash.Hash
			messageHashes[transition.RespondingMessage] = respondingMessageHash
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

func validateTransitionExecutable(instance Instance, transition Transition, input ConstraintInput) error {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if instance.TokenCounts[incomingPlaceId] < 1 {
			return fmt.Errorf("transition %s is not executable because there are not enough tokens", transition.Id)
		}
	}
	return validateConstraint(instance, transition, input)
}

func (instance *Instance) FindMessageHashById(id ModelMessageId) Hash {
	return instance.MessageHashes[id]
}

func (instance *Instance) FindPublicKeyByParticipant(id ParticipantId) PublicKey {
	return instance.PublicKeys[id]
}

func validateConstraint(instance Instance, transition Transition, input ConstraintInput) error {
	constraint := transition.Constraint
	if len(constraint.MessageIds) != len(input.Messages) {
		return fmt.Errorf("transition %s is not executable because number of constraint inputs differs from the number of messages in the constraint", transition.Id)
	}
	lhs := constraint.Offset
	for i, message := range input.Messages {
		hash := message.Hash.Hash
		messageId := EmptyMessageId
		for i, messageHash := range instance.MessageHashes {
			if bytes.Equal(hash.Value[:], messageHash.Value[:]) {
				messageId = ModelMessageId(i)
				break
			}
		}
		if constraint.Coefficients[i] != 0 && messageId != constraint.MessageIds[i] {
			return fmt.Errorf("transition %s is not executable because the wrong constraint inputs have been provided", transition.Id)
		}
		lhs += constraint.Coefficients[i] * input.Messages[i].IntegerMessage
	}

	var result bool
	switch constraint.ComparisonOperator {
	case OperatorEqual:
		result = lhs == 0
	case OperatorGreaterThan:
		result = lhs > 0
	case OperatorLessThan:
		result = lhs < 0
	case OperatorGreaterThanOrEqual:
		result = lhs >= 0
	case OperatorLessThanOrEqual:
		result = lhs <= 0
	}
	if !result {
		return fmt.Errorf("transition %s is not executable because the constraint evaluates to false", transition.Id)
	}
	return nil
}
