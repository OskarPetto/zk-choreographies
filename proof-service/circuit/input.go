package circuit

import (
	"fmt"
	"proof-service/crypto"
	"proof-service/workflow"

	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type Commitment struct {
	Value      frontend.Variable `gnark:",public"`
	Randomness frontend.Variable
}

type Signature struct {
	Value     eddsa.Signature
	PublicKey eddsa.PublicKey
}

type Participant struct {
	PublicKey eddsa.PublicKey
}

type Instance struct {
	PlaceCount       frontend.Variable `gnark:",public"`
	TokenCounts      [workflow.MaxPlaceCount]frontend.Variable
	ParticipantCount frontend.Variable `gnark:",public"`
	Participants     [workflow.MaxParticipantCount]Participant
}

type Transition struct {
	RequiresParticipant frontend.Variable                              `gnark:",public"`
	Participant         frontend.Variable                              `gnark:",public"`
	IncomingPlaceCount  frontend.Variable                              `gnark:",public"`
	IncomingPlaces      [workflow.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaceCount  frontend.Variable                              `gnark:",public"`
	OutgoingPlaces      [workflow.MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount       frontend.Variable `gnark:",public"`
	StartPlace       frontend.Variable `gnark:",public"`
	EndPlace         frontend.Variable `gnark:",public"`
	TransitionCount  frontend.Variable `gnark:",public"`
	Transitions      [workflow.MaxTransitionCount]Transition
	ParticipantCount frontend.Variable `gnark:",public"`
}

func FromCommitment(commitment crypto.Commitment) Commitment {
	return Commitment{
		Value:      commitment.Value,
		Randomness: commitment.Randomness,
	}
}

func FromSignature(signature crypto.Signature) Signature {
	var value eddsa.Signature
	value.Assign(twistededwards.BN254, signature.Value)
	var publicKey eddsa.PublicKey
	publicKey.Assign(twistededwards.BN254, signature.PublicKey)
	return Signature{
		Value:     value,
		PublicKey: publicKey,
	}
}

func FromInstance(instance workflow.Instance) (Instance, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > workflow.MaxPlaceCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many places", instance.Id)
	}
	var tokenCounts [workflow.MaxPlaceCount]frontend.Variable
	for i := 0; i < placeCount; i++ {
		tokenCounts[i] = instance.TokenCounts[i]
	}
	for i := placeCount; i < workflow.MaxPlaceCount; i++ {
		tokenCounts[i] = 0
	}
	participantCount := len(instance.Participants)
	if participantCount > workflow.MaxParticipantCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many participants", instance.Id)
	}
	var participants [workflow.MaxParticipantCount]Participant
	for i := 0; i < participantCount; i++ {
		participants[i] = fromParticipant(instance.Participants[i])
	}
	for i := participantCount; i < workflow.MaxParticipantCount; i++ {
		participants[i] = emptyParticipant()
	}
	return Instance{
		PlaceCount:       placeCount,
		TokenCounts:      tokenCounts,
		ParticipantCount: participantCount,
		Participants:     participants,
	}, nil
}

func fromParticipant(participant workflow.Participant) Participant {
	var publicKey eddsa.PublicKey
	publicKey.Assign(twistededwards.BN254, participant.PublicKey)
	return Participant{
		PublicKey: publicKey,
	}
}

func emptyParticipant() Participant {
	var publicKey eddsa.PublicKey
	zeros := make([]byte, 32)
	publicKey.Assign(twistededwards.BN254, zeros)
	return Participant{
		PublicKey: publicKey,
	}
}

func FromPetriNet(petriNet workflow.PetriNet) (PetriNet, error) {
	placeCount := petriNet.PlaceCount
	if placeCount > workflow.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has too many places", petriNet.Id)
	}
	transitions, err := fromTransitions(petriNet.Id, petriNet.Transitions)
	if err != nil {
		return PetriNet{}, err
	}
	if petriNet.StartPlace >= workflow.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has invalid startPlace", petriNet.Id)
	}
	if petriNet.EndPlace >= workflow.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has invalid endPlace", petriNet.Id)
	}
	return PetriNet{
		PlaceCount:       petriNet.PlaceCount,
		StartPlace:       petriNet.StartPlace,
		EndPlace:         petriNet.EndPlace,
		TransitionCount:  len(petriNet.Transitions),
		Transitions:      transitions,
		ParticipantCount: petriNet.ParticipantCount,
	}, nil
}

func fromTransitions(petriNetId string, workflowTransitions []workflow.Transition) ([workflow.MaxTransitionCount]Transition, error) {
	transitionCount := len(workflowTransitions)
	if transitionCount > workflow.MaxTransitionCount {
		return [workflow.MaxTransitionCount]Transition{}, fmt.Errorf("petriNet '%s' has too many transitions", petriNetId)
	}
	var transitions [workflow.MaxTransitionCount]Transition
	var err error
	for i := 0; i < transitionCount; i++ {
		transitions[i], err = fromTransition(workflowTransitions[i])
		if err != nil {
			return [workflow.MaxTransitionCount]Transition{}, fmt.Errorf("petriNet '%s' cannot be mapped because transition at index %d is invalid: %w", petriNetId, i, err)
		}
	}
	for i := transitionCount; i < workflow.MaxTransitionCount; i++ {
		transitions[i] = emptyTransition()
	}
	return transitions, nil
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
		incomingPlaces[i] = transition.IncomingPlaces[i]
	}
	for i := incomingPlaceCount; i < workflow.MaxBranchingFactor; i++ {
		incomingPlaces[i] = workflow.MaxPlaceCount
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		outgoingPlaces[i] = transition.OutgoingPlaces[i]
	}
	for i := outgoingPlaceCount; i < workflow.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = workflow.MaxPlaceCount
	}
	requiresParticipant := 0
	if transition.RequiresParticipant {
		requiresParticipant = 1
	}
	return Transition{
		IncomingPlaceCount:  incomingPlaceCount,
		IncomingPlaces:      incomingPlaces,
		OutgoingPlaceCount:  outgoingPlaceCount,
		OutgoingPlaces:      outgoingPlaces,
		RequiresParticipant: requiresParticipant,
		Participant:         transition.Participant,
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
		IncomingPlaceCount:  0,
		IncomingPlaces:      incomingPlaces,
		OutgoingPlaceCount:  0,
		OutgoingPlaces:      outgoingPlaces,
		RequiresParticipant: 0,
		Participant:         0,
	}
}
