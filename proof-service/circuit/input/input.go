package input

import (
	"proof-service/commitment"
	"proof-service/instance"
	"proof-service/petri_net"

	"github.com/consensys/gnark/std/math/uints"
)

type Commitment struct {
	Value      [commitment.CommitmentSize]uints.U8 `gnark:",public"`
	Randomness []uints.U8
}

type Instance struct {
	PlaceCount  uints.U8
	TokenCounts []uints.U8
}

type Transition struct {
	IncomingPlaceCount uints.U8   `gnark:",public"`
	IncomingPlaces     []uints.U8 `gnark:",public"`
	OutgoingPlaceCount uints.U8   `gnark:",public"`
	OutgoingPlaces     []uints.U8 `gnark:",public"`
}

type PetriNet struct {
	PlaceCount      uints.U8 `gnark:",public"`
	StartPlace      uints.U8 `gnark:",public"`
	TransitionCount uints.U8 `gnark:",public"`
	Transitions     []Transition
}

func FromCommitment(c commitment.Commitment) Commitment {
	return Commitment{
		Value:      ([commitment.CommitmentSize]uints.U8)(uints.NewU8Array(c.Value[:])),
		Randomness: uints.NewU8Array(c.Randomness),
	}
}

func FromInstance(inst instance.Instance) Instance {
	tokenCounts := make([]uints.U8, inst.PlaceCount)
	for i := 0; i < int(inst.PlaceCount); i++ {
		tokenCounts[i] = uints.NewU8(byte(inst.TokenCounts[i]))
	}
	return Instance{
		PlaceCount:  uints.NewU8(inst.PlaceCount),
		TokenCounts: tokenCounts,
	}
}

func FromPetriNet(petriNet petri_net.PetriNet) PetriNet {
	transitions := make([]Transition, petriNet.TransitionCount)
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

func fromTransition(transition petri_net.Transition) Transition {
	incomingPlaces := uints.NewU8Array(transition.IncomingPlaces[:])
	outgoingPlaces := uints.NewU8Array(transition.OutgoingPlaces[:])
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
