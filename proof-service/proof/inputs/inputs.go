package inputs

import (
	"fmt"
	"proof-service/commitment"
	"proof-service/workflow"

	"github.com/consensys/gnark/std/math/uints"
)

const MaxPlaceCount = 100
const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type Commitment struct {
	Value      [commitment.CommitmentSize]uints.U8 `gnark:",public"`
	Randomness [commitment.RandomnessSize]uints.U8
}

type Instance struct {
	PlaceCount  uints.U8
	TokenCounts [MaxPlaceCount]uints.U8
}

type Transition struct {
	IncomingPlaceCount uints.U8                     `gnark:",public"`
	IncomingPlaces     [MaxBranchingFactor]uints.U8 `gnark:",public"`
	OutgoingPlaceCount uints.U8                     `gnark:",public"`
	OutgoingPlaces     [MaxBranchingFactor]uints.U8 `gnark:",public"`
}

type PetriNet struct {
	PlaceCount      uints.U8 `gnark:",public"`
	StartPlace      uints.U8 `gnark:",public"`
	TransitionCount uints.U8 `gnark:",public"`
	Transitions     [MaxTransitionCount]Transition
}

func FromCommitment(c commitment.Commitment) (Commitment, error) {
	randomnessSize := len(c.Randomness)
	if randomnessSize > commitment.RandomnessSize {
		return Commitment{}, fmt.Errorf("commitment %s is too large", c.Id)
	}
	var randomness [commitment.RandomnessSize]uints.U8
	for i := 0; i < commitment.RandomnessSize; i++ {
		randomness[i] = uints.NewU8(c.Randomness[i])
	}
	var value [commitment.CommitmentSize]uints.U8
	for i := 0; i < commitment.CommitmentSize; i++ {
		value[i] = uints.NewU8(c.Value[i])
	}
	return Commitment{
		Value:      value,
		Randomness: randomness,
	}, nil
}

func FromInstance(instance workflow.Instance) (Instance, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > MaxPlaceCount {
		return Instance{}, fmt.Errorf("instance %s is too large", instance.Id)
	}
	var tokenCounts [MaxPlaceCount]uints.U8
	for i := 0; i < placeCount; i++ {
		tokenCounts[i] = uints.NewU8(uint8(instance.TokenCounts[i]))
	}
	return Instance{
		PlaceCount:  uints.NewU8(uint8(placeCount)),
		TokenCounts: tokenCounts,
	}, nil
}

func FromPetriNet(petriNet workflow.PetriNet) (PetriNet, error) {
	placeCount := petriNet.PlaceCount
	transitionCount := len(petriNet.Transitions)
	if placeCount > MaxPlaceCount || transitionCount > MaxTransitionCount {
		return PetriNet{}, fmt.Errorf("petriNet %s is too large", petriNet.Id)
	}
	var transitions [MaxTransitionCount]Transition
	var err error
	for i := 0; i < transitionCount; i++ {
		transitions[i], err = fromTransition(petriNet.Transitions[i])
		if err != nil {
			return PetriNet{}, fmt.Errorf("petriNet %s cannot be mapped because: %w", petriNet.Id, err)
		}
	}
	return PetriNet{
		PlaceCount:      uints.NewU8(uint8(placeCount)),
		StartPlace:      uints.NewU8(petriNet.StartPlace),
		TransitionCount: uints.NewU8(uint8(transitionCount)),
		Transitions:     transitions,
	}, nil
}

func fromTransition(transition workflow.Transition) (Transition, error) {
	incomingPlaceCount := len(transition.IncomingPlaces)
	outgoingPlaceCount := len(transition.OutgoingPlaces)
	if incomingPlaceCount > MaxBranchingFactor || outgoingPlaceCount > MaxBranchingFactor {
		return Transition{}, fmt.Errorf("transition %s branches to much", transition.Id)
	}
	var incomingPlaces [MaxBranchingFactor]uints.U8
	var outgoingPlaces [MaxBranchingFactor]uints.U8
	for i := 0; i < incomingPlaceCount; i++ {
		incomingPlaces[i] = uints.NewU8(transition.IncomingPlaces[i])
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		outgoingPlaces[i] = uints.NewU8(transition.OutgoingPlaces[i])
	}
	return Transition{
		IncomingPlaceCount: uints.NewU8(uint8(incomingPlaceCount)),
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: uints.NewU8(uint8(outgoingPlaceCount)),
		OutgoingPlaces:     outgoingPlaces,
	}, nil
}
