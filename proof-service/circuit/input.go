package circuit

import (
	"proof-service/commitment"
	"proof-service/domain"

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
	IncomingPlaceCount uints.U8                            `gnark:",public"`
	IncomingPlaces     [domain.MaxBranchingFactor]uints.U8 `gnark:",public"`
	OutgoingPlaceCount uints.U8                            `gnark:",public"`
	OutgoingPlaces     [domain.MaxBranchingFactor]uints.U8 `gnark:",public"`
}

type PetriNet struct {
	PlaceCount      uints.U8 `gnark:",public"`
	StartPlace      uints.U8 `gnark:",public"`
	TransitionCount uints.U8 `gnark:",public"`
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
	return PetriNet{
		PlaceCount:      uints.NewU8(petriNet.PlaceCount),
		StartPlace:      uints.NewU8(petriNet.StartPlace),
		TransitionCount: uints.NewU8(petriNet.TransitionCount),
		Transitions:     transitions,
	}
}

func fromTransition(transition domain.Transition) Transition {
	var incomingPlaces [domain.MaxBranchingFactor]uints.U8
	var outgoingPlaces [domain.MaxBranchingFactor]uints.U8
	for i := 0; i < int(transition.IncomingPlaceCount); i++ {
		incomingPlaces[i] = uints.NewU8(transition.IncomingPlaces[i])
	}
	for i := 0; i < int(transition.OutgoingPlaceCount); i++ {
		outgoingPlaces[i] = uints.NewU8(transition.OutgoingPlaces[i])
	}
	return Transition{
		IncomingPlaceCount: uints.NewU8(transition.IncomingPlaceCount),
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: uints.NewU8(transition.OutgoingPlaceCount),
		OutgoingPlaces:     outgoingPlaces,
	}
}
