package petri_net

import (
	"fmt"
)

const MaxPlaceCount = 100
const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type PlaceId = uint8

type Transition struct {
	Id             string
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type PetriNet struct {
	Id          string
	PlaceCount  uint8
	StartPlace  PlaceId
	Transitions []Transition
}

func ValidatePetriNet(petriNet PetriNet) error {
	placeCount := petriNet.PlaceCount
	transitionCount := len(petriNet.Transitions)
	if placeCount > MaxPlaceCount || transitionCount > MaxTransitionCount {
		return fmt.Errorf("petriNet '%s' is too large", petriNet.Id)
	}
	if petriNet.StartPlace >= MaxPlaceCount {
		return fmt.Errorf("petriNet '%s' has invalid startPlace", petriNet.Id)
	}
	for i := 0; i < transitionCount; i++ {
		err := ValidateTransition(petriNet.Transitions[i])
		if err != nil {
			return fmt.Errorf("petriNet '%s' cannot be mapped because transition at index %d is invalid: %w", petriNet.Id, i, err)
		}
	}
	return nil
}

func ValidateTransition(transition Transition) error {
	incomingPlaceCount := len(transition.IncomingPlaces)
	outgoingPlaceCount := len(transition.OutgoingPlaces)
	if incomingPlaceCount > MaxBranchingFactor || outgoingPlaceCount > MaxBranchingFactor {
		return fmt.Errorf("transition '%s' branches too much", transition.Id)
	}
	for i := 0; i < incomingPlaceCount; i++ {
		if transition.IncomingPlaces[i] >= MaxPlaceCount {
			return fmt.Errorf("incomingPlace at index %d of transition '%s' is invalid", i, transition.Id)
		}
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		if transition.OutgoingPlaces[i] >= MaxPlaceCount {
			return fmt.Errorf("outgoingPlace at index %d of transition '%s' is invalid", i, transition.Id)
		}
	}
	return nil
}
