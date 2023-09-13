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
var executionService = execution.ExecutionService{}

func TestInstantiateModel(t *testing.T) {
	result := executionService.InstantiateModel(model1)
	assert.Equal(t, instance1.Model, result.Model)
	assert.Equal(t, instance1.TokenCounts, result.TokenCounts)
}

func TestImmutabilityOfInstances(t *testing.T) {
	trace := []model.TransitionId{"As"}
	transitions := findTransitions(model1, trace)
	executionService.ExecuteTransition(instance1, transitions[0])
	assert.Equal(t, []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, instance1.TokenCounts)
}

func TestExecuteTransition(t *testing.T) {
	trace := []model.TransitionId{"As", "Aa", "Fa", "Sso", "Ro", "Ao", "Aaa", "Af"}
	transitions := findTransitions(model1, trace)
	instance := instance1
	for _, transition := range transitions {
		instance, _ = executionService.ExecuteTransition(instance, transition)
	}
	assert.Equal(t, []int8{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, instance.TokenCounts)
}
