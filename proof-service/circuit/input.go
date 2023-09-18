package circuit

import (
	"proof-service/commitment"
	"proof-service/domain"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
)

type Commitment struct {
	Value      [commitment.CommitmentSize]uints.U8 `gnark:",public"`
	Randomness [commitment.RandomnessSize]uints.U8
}

type Instance struct {
	PlaceCount  uints.U8 `gnark:",public"`
	TokenCounts [domain.MaxPlaceCount]uints.U8
}

type Transition struct {
	IncomingPlaceCount frontend.Variable                            `gnark:",public"`
	IncomingPlaces     [domain.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaceCount frontend.Variable                            `gnark:",public"`
	OutgoingPlaces     [domain.MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount      frontend.Variable `gnark:",public"`
	StartPlace      frontend.Variable `gnark:",public"`
	TransitionCount frontend.Variable `gnark:",public"`
	Transitions     [domain.MaxTransitionCount]Transition
}

func FromCommitment(c commitment.Commitment) Commitment {
	return Commitment{
		Value:      ([commitment.CommitmentSize]uints.U8)(uints.NewU8Array(c.Value[:])),
		Randomness: ([commitment.RandomnessSize]uints.U8)(uints.NewU8Array(c.Randomness[:])),
	}
}

func FromInstance(inst domain.Instance) Instance {
	var tokenCounts [domain.MaxPlaceCount]uints.U8
	for i := 0; i < int(inst.PlaceCount); i++ {
		tokenCounts[i] = uints.NewU8(byte(inst.TokenCounts[i]))
	}
	return Instance{
		PlaceCount:  uints.NewU8(inst.PlaceCount),
		TokenCounts: tokenCounts,
	}
}

func FromPetriNet(petriNet domain.PetriNet) PetriNet {
	var transitions [domain.MaxTransitionCount]Transition
	for i := 0; i < int(petriNet.TransitionCount); i++ {
		transitions[i] = fromTransition(petriNet.Transitions[i])
	}
	for i := int(petriNet.TransitionCount); i < domain.MaxTransitionCount; i++ {
		transitions[i] = emptyTransition()
	}
	return PetriNet{
		PlaceCount:      petriNet.PlaceCount,
		StartPlace:      petriNet.StartPlace,
		TransitionCount: petriNet.TransitionCount,
		Transitions:     transitions,
	}
}

func fromTransition(transition domain.Transition) Transition {
	var incomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i := 0; i < int(transition.IncomingPlaceCount); i++ {
		incomingPlaces[i] = transition.IncomingPlaces[i]
	}
	for i := int(transition.IncomingPlaceCount); i < domain.MaxBranchingFactor; i++ {
		incomingPlaces[i] = 0
	}
	for i := 0; i < int(transition.OutgoingPlaceCount); i++ {
		outgoingPlaces[i] = transition.OutgoingPlaces[i]
	}
	for i := int(transition.OutgoingPlaceCount); i < domain.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = 0
	}
	return Transition{
		IncomingPlaceCount: transition.IncomingPlaceCount,
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: transition.OutgoingPlaceCount,
		OutgoingPlaces:     outgoingPlaces,
	}
}

func emptyTransition() Transition {
	var incomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i := 0; i < domain.MaxBranchingFactor; i++ {
		incomingPlaces[i] = 0
	}
	for i := 0; i < domain.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = 0
	}
	return Transition{
		IncomingPlaceCount: 0,
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: 0,
		OutgoingPlaces:     outgoingPlaces,
	}
}
