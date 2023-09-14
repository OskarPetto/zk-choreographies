package petri_net

import (
	"fmt"
	"proof-service/workflow"
)

const MaxPlaceCount = 100
const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type PlaceId = uint8

type Transition struct {
	Id                 string
	IncomingPlaceCount uint8
	IncomingPlaces     [MaxBranchingFactor]PlaceId
	OutgoingPlaceCount uint8
	OutgoingPlaces     [MaxBranchingFactor]PlaceId
}

type PetriNet struct {
	Id              string
	PlaceCount      uint8
	StartPlace      PlaceId
	TransitionCount uint8
	Transitions     [MaxTransitionCount]Transition
}

func FromWorkflowPetriNet(petriNet workflow.PetriNet) (PetriNet, error) {
	placeCount := petriNet.PlaceCount
	transitionCount := len(petriNet.Transitions)
	if placeCount > MaxPlaceCount || transitionCount > MaxTransitionCount {
		return PetriNet{}, fmt.Errorf("petriNet %s is too large", petriNet.Id)
	}
	if petriNet.StartPlace >= MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet %s has invalid startPlace", petriNet.Id)
	}
	var transitions [MaxTransitionCount]Transition
	var err error
	for i := 0; i < transitionCount; i++ {
		transitions[i], err = fromWorkflowTransition(petriNet.Transitions[i])
		if err != nil {
			return PetriNet{}, fmt.Errorf("petriNet %s cannot be mapped because transition at index %d is invalid: %w", petriNet.Id, i, err)
		}
	}
	return PetriNet{
		PlaceCount:      uint8(placeCount),
		StartPlace:      uint8(petriNet.StartPlace),
		TransitionCount: uint8(transitionCount),
		Transitions:     transitions,
	}, nil
}

func fromWorkflowTransition(transition workflow.Transition) (Transition, error) {
	incomingPlaceCount := len(transition.IncomingPlaces)
	outgoingPlaceCount := len(transition.OutgoingPlaces)
	if incomingPlaceCount > MaxBranchingFactor || outgoingPlaceCount > MaxBranchingFactor {
		return Transition{}, fmt.Errorf("transition %s branches too much", transition.Id)
	}
	var incomingPlaces [MaxBranchingFactor]uint8
	var outgoingPlaces [MaxBranchingFactor]uint8
	for i := 0; i < incomingPlaceCount; i++ {
		if transition.IncomingPlaces[i] >= MaxBranchingFactor {
			return Transition{}, fmt.Errorf("incomingPlace at index %d of transition %s is invalid", i, transition.Id)
		}
		incomingPlaces[i] = uint8(transition.IncomingPlaces[i])
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		if transition.OutgoingPlaces[i] >= MaxBranchingFactor {
			return Transition{}, fmt.Errorf("outgoingPlace at index %d of transition %s is invalid", i, transition.Id)
		}
		outgoingPlaces[i] = uint8(transition.OutgoingPlaces[i])
	}
	return Transition{
		IncomingPlaceCount: uint8(incomingPlaceCount),
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: uint8(outgoingPlaceCount),
		OutgoingPlaces:     outgoingPlaces,
	}, nil
}
