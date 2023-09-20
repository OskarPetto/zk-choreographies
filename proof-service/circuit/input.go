package circuit

import (
	"fmt"
	"proof-service/crypto"
	"proof-service/workflow"

	"github.com/consensys/gnark/frontend"
)

type Commitment struct {
	Value      frontend.Variable `gnark:",public"`
	Randomness frontend.Variable
}

type Instance struct {
	PlaceCount  frontend.Variable `gnark:",public"`
	TokenCounts [workflow.MaxPlaceCount]frontend.Variable
}

type Transition struct {
	IncomingPlaceCount frontend.Variable                              `gnark:",public"`
	IncomingPlaces     [workflow.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaceCount frontend.Variable                              `gnark:",public"`
	OutgoingPlaces     [workflow.MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount      frontend.Variable `gnark:",public"`
	StartPlace      frontend.Variable `gnark:",public"`
	EndPlace        frontend.Variable `gnark:",public"`
	TransitionCount frontend.Variable `gnark:",public"`
	Transitions     [workflow.MaxTransitionCount]Transition
}

func FromCommitment(commitment crypto.Commitment) Commitment {
	return Commitment{
		Value:      commitment.Value,
		Randomness: commitment.Randomness,
	}
}

func FromInstance(instance workflow.Instance) (Instance, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > workflow.MaxPlaceCount {
		return Instance{}, fmt.Errorf("instance '%s' is too large", instance.Id)
	}
	var tokenCounts [workflow.MaxPlaceCount]frontend.Variable
	for i := 0; i < placeCount; i++ {
		tokenCounts[i] = instance.TokenCounts[i]
	}
	for i := placeCount; i < workflow.MaxPlaceCount; i++ {
		tokenCounts[i] = 0
	}
	return Instance{
		PlaceCount:  placeCount,
		TokenCounts: tokenCounts,
	}, nil
}

func FromPetriNet(petriNet workflow.PetriNet) (PetriNet, error) {
	placeCount := petriNet.PlaceCount
	transitionCount := len(petriNet.Transitions)
	if placeCount > workflow.MaxPlaceCount || transitionCount > workflow.MaxTransitionCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' is too large", petriNet.Id)
	}
	if petriNet.StartPlace >= workflow.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has invalid startPlace", petriNet.Id)
	}
	var transitions [workflow.MaxTransitionCount]Transition
	var err error
	for i := 0; i < transitionCount; i++ {
		transitions[i], err = fromTransition(petriNet.Transitions[i])
		if err != nil {
			return PetriNet{}, fmt.Errorf("petriNet '%s' cannot be mapped because transition at index %d is invalid: %w", petriNet.Id, i, err)
		}
	}
	for i := transitionCount; i < workflow.MaxTransitionCount; i++ {
		transitions[i] = emptyTransition()
	}
	return PetriNet{
		PlaceCount:      petriNet.PlaceCount,
		StartPlace:      petriNet.StartPlace,
		EndPlace:        petriNet.EndPlace,
		TransitionCount: transitionCount,
		Transitions:     transitions,
	}, nil
}

func fromTransition(transition workflow.Transition) (Transition, error) {
	incomingPlaceCount := len(transition.IncomingPlaces)
	outgoingPlaceCount := len(transition.OutgoingPlaces)
	if incomingPlaceCount > workflow.MaxBranchingFactor || outgoingPlaceCount > workflow.MaxBranchingFactor {
		return Transition{}, fmt.Errorf("transition '%s' branches too much", transition.Id)
	}
	var incomingPlaces [workflow.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [workflow.MaxBranchingFactor]frontend.Variable
	for i := 0; i < incomingPlaceCount; i++ {
		if transition.IncomingPlaces[i] >= workflow.MaxPlaceCount {
			return Transition{}, fmt.Errorf("incomingPlace at index %d of transition '%s' is invalid", i, transition.Id)
		}
		incomingPlaces[i] = transition.IncomingPlaces[i]
	}
	for i := incomingPlaceCount; i < workflow.MaxBranchingFactor; i++ {
		incomingPlaces[i] = workflow.MaxPlaceCount
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		if transition.OutgoingPlaces[i] >= workflow.MaxPlaceCount {
			return Transition{}, fmt.Errorf("outgoingPlace at index %d of transition '%s' is invalid", i, transition.Id)
		}
		outgoingPlaces[i] = transition.OutgoingPlaces[i]
	}
	for i := outgoingPlaceCount; i < workflow.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = workflow.MaxPlaceCount
	}
	return Transition{
		IncomingPlaceCount: incomingPlaceCount,
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: outgoingPlaceCount,
		OutgoingPlaces:     outgoingPlaces,
	}, nil
}

func emptyTransition() Transition {
	var incomingPlaces [workflow.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [workflow.MaxBranchingFactor]frontend.Variable
	for i := 0; i < workflow.MaxBranchingFactor; i++ {
		incomingPlaces[i] = workflow.MaxPlaceCount
	}
	for i := 0; i < workflow.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = workflow.MaxPlaceCount
	}
	return Transition{
		IncomingPlaceCount: 0,
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: 0,
		OutgoingPlaces:     outgoingPlaces,
	}
}
