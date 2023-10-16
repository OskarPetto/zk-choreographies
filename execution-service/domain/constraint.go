package domain

import (
	"bytes"
	"math"
)

type IntegerType = int32

const EmptyIntegerMessage = IntegerType(math.MinInt32)
const MaxConstraintMessageCount = 2

type ComparisonOperator = uint8

const (
	OperatorEqual              = 0
	OperatorGreaterThan        = 1
	OperatorLessThan           = 2
	OperatorGreaterThanOrEqual = 3
	OperatorLessThanOrEqual    = 4
)

var ValidComparisonOperators = []ComparisonOperator{OperatorEqual, OperatorGreaterThan, OperatorLessThan, OperatorGreaterThanOrEqual, OperatorLessThanOrEqual}

// ax + by + c = 0
type Constraint struct {
	Coefficients       []IntegerType
	MessageIds         []ModelMessageId
	Offset             IntegerType
	ComparisonOperator ComparisonOperator
}

func EmptyConstraint() Constraint {
	return Constraint{}
}

type ConstraintInput struct {
	Messages []Message
}

func EmptyConstraintInput() ConstraintInput {
	return ConstraintInput{
		Messages: make([]Message, 0),
	}
}

func (instance *Instance) EvaluateConstraint(constraint Constraint, input ConstraintInput) bool {
	if len(constraint.MessageIds) != len(input.Messages) {
		return false
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
			return false
		}
		lhs += constraint.Coefficients[i] * input.Messages[i].IntegerMessage
	}

	switch constraint.ComparisonOperator {
	case OperatorEqual:
		return lhs == 0
	case OperatorGreaterThan:
		return lhs > 0
	case OperatorLessThan:
		return lhs < 0
	case OperatorGreaterThanOrEqual:
		return lhs >= 0
	case OperatorLessThanOrEqual:
		return lhs <= 0
	}
	return false
}
