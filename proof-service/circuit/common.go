package circuit

import (
	"proof-service/commitment"
	"proof-service/instance"
	"proof-service/petri_net"

	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha2"
	"github.com/consensys/gnark/std/math/uints"
)

type Commitment struct {
	Value      [commitment.CommitmentSize]uints.U8 `gnark:",public"`
	Randomness [commitment.RandomnessSize]uints.U8
}

type Instance struct {
	PlaceCount  uints.U8 `gnark:",public"`
	TokenCounts [petri_net.MaxPlaceCount]uints.U8
}

type Transition struct {
	IncomingPlaceCount uints.U8                               `gnark:",public"`
	IncomingPlaces     [petri_net.MaxBranchingFactor]uints.U8 `gnark:",public"`
	OutgoingPlaceCount uints.U8                               `gnark:",public"`
	OutgoingPlaces     [petri_net.MaxBranchingFactor]uints.U8 `gnark:",public"`
}

type PetriNet struct {
	PlaceCount      uints.U8 `gnark:",public"`
	StartPlace      uints.U8 `gnark:",public"`
	TransitionCount uints.U8 `gnark:",public"`
	Transitions     [petri_net.MaxTransitionCount]Transition
}

func FromCommitment(c commitment.Commitment) Commitment {
	return Commitment{
		Value:      ([commitment.CommitmentSize]uints.U8)(uints.NewU8Array(c.Value[:])),
		Randomness: ([commitment.RandomnessSize]uints.U8)(uints.NewU8Array(c.Randomness[:])),
	}
}

func FromInstance(inst instance.Instance) Instance {
	var tokenCounts [petri_net.MaxPlaceCount]uints.U8
	for i := 0; i < int(inst.PlaceCount); i++ {
		tokenCounts[i] = uints.NewU8(byte(inst.TokenCounts[i]))
	}
	return Instance{
		PlaceCount:  uints.NewU8(inst.PlaceCount),
		TokenCounts: tokenCounts,
	}
}

func FromPetriNet(petriNet petri_net.PetriNet) PetriNet {
	var transitions [petri_net.MaxTransitionCount]Transition
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
	var incomingPlaces [petri_net.MaxBranchingFactor]uints.U8
	var outgoingPlaces [petri_net.MaxBranchingFactor]uints.U8
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

func checkCommitment(api frontend.API, uapi *uints.BinaryField[uints.U32], instance Instance, commitment Commitment) error {
	hasher, err := sha2.New(api)
	if err != nil {
		return err
	}
	hasher.Write([]uints.U8{instance.PlaceCount})
	hasher.Write(instance.TokenCounts[:])
	hasher.Write(commitment.Randomness[:])
	hash := hasher.Sum()
	for i := range commitment.Value {
		uapi.ByteAssertEq(commitment.Value[i], hash[i])
	}
	return nil
}

func equals(api frontend.API, a, b frontend.Variable) frontend.Variable {
	return api.IsZero(api.Sub(a, b))
}
