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
	Coefficients       [MaxConstraintMessageCount]IntegerType
	MessageIds         [MaxConstraintMessageCount]MessageId
	Offset             IntegerType
	ComparisonOperator ComparisonOperator
}

func EmptyConstraint() Constraint {
	return Constraint{
		MessageIds: [MaxConstraintMessageCount]MessageId{EmptyMessageId, EmptyMessageId},
	}
}

type ConstraintInput struct {
	IntegerMessages [MaxConstraintMessageCount]Message
}

func EmptyConstraintInput() ConstraintInput {
	return ConstraintInput{}
}

func (instance *Instance) EvaluateConstraint(constraint Constraint, input ConstraintInput) bool {
	lhs := constraint.Offset
	for i := 0; i < MaxConstraintMessageCount; i++ {
		hash := input.IntegerMessages[i].Hash.Value
		messageId := EmptyMessageId
		for i, messageHash := range instance.MessageHashes {
			if bytes.Equal(hash[:], messageHash[:]) {
				messageId = MessageId(i)
				break
			}
		}
		if messageId != constraint.MessageIds[i] {
			return false
		}
		lhs += constraint.Coefficients[i] * input.IntegerMessages[i].IntegerMessage
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
