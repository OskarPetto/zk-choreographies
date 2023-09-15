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
	TokenCounts []uints.U8
}

type Transition struct {
	IncomingPlaces []uints.U8 `gnark:",public"`
	OutgoingPlaces []uints.U8 `gnark:",public"`
}

type PetriNet struct {
	PlaceCount  uints.U8 `gnark:",public"`
	StartPlace  uints.U8 `gnark:",public"`
	Transitions []Transition
}

func FromCommitment(c commitment.Commitment) Commitment {
	return Commitment{
		Value:      ([commitment.CommitmentSize]uints.U8)(uints.NewU8Array(c.Value[:])),
		Randomness: uints.NewU8Array(c.Randomness),
	}
}

func FromInstance(inst instance.Instance) Instance {
	tokenCounts := make([]uints.U8, len(inst.TokenCounts))
	for i := 0; i < len(inst.TokenCounts); i++ {
		tokenCounts[i] = uints.NewU8(byte(inst.TokenCounts[i]))
	}
	return Instance{
		TokenCounts: tokenCounts,
	}
}

func FromPetriNet(petriNet petri_net.PetriNet) PetriNet {
	transitions := make([]Transition, len(petriNet.Transitions))
	for i := 0; i < int(len(petriNet.Transitions)); i++ {
		transitions[i] = fromTransition(petriNet.Transitions[i])
	}
	return PetriNet{
		PlaceCount:  uints.NewU8(petriNet.PlaceCount),
		StartPlace:  uints.NewU8(petriNet.StartPlace),
		Transitions: transitions,
	}
}

func fromTransition(transition petri_net.Transition) Transition {
	return Transition{
		IncomingPlaces: uints.NewU8Array(transition.IncomingPlaces[:]),
		OutgoingPlaces: uints.NewU8Array(transition.OutgoingPlaces[:]),
	}
}
