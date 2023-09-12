package execution_test

import (
	"execution-service/execution"
	"execution-service/model"
	mdl "execution-service/model"
	"execution-service/testdata"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func findTransitions(model mdl.Model, transitionIds []mdl.TransitionId) []mdl.Transition {
	transitions := make([]mdl.Transition, len(transitionIds))
	for i := range transitionIds {
		transitionId := transitionIds[i]
		found := false
		for _, transition := range model.Transitions {
			if transition.Id == transitionId {
				transitions[i] = transition
				found = true
				break
			}
		}
		if !found {
			panic(fmt.Sprintf("Transition %s does not exist in model %+v", transitionId, model))
		}
	}
	return transitions
}

var model1 = testdata.GetModel1()
var instance1 = testdata.GetInstance1()

func TestInstantiateModel(t *testing.T) {
	result := execution.InstantiateModel(model1)
	assert.Equal(t, instance1.Model, result.Model)
	assert.Equal(t, instance1.TokenCounts, result.TokenCounts)
}

func TestExecuteTransition(t *testing.T) {
	trace := []model.TransitionId{"As", "Aa", "Fa", "Sso", "Ro", "Ao", "Aaa", "Af"}
	transitions := findTransitions(model1, trace)
	instance := instance1
	for _, transition := range transitions {
		instance, _ = execution.ExecuteTransition(instance, transition)
	}
	assert.Equal(t, instance.TokenCounts, []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
