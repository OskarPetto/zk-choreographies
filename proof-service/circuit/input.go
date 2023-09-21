package circuit

import (
	"fmt"
	"proof-service/crypto"
	"proof-service/workflow"

	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type SaltedHash struct {
	Value frontend.Variable `gnark:",public"`
	Salt  frontend.Variable
}

type Signature struct {
	Value     eddsa.Signature
	PublicKey eddsa.PublicKey
}

type Instance struct {
	TokenCounts [workflow.MaxPlaceCount]frontend.Variable
	PublicKeys  [workflow.MaxParticipantCount]eddsa.PublicKey
}

type Transition struct {
	IsExecutableByAnyParticipant frontend.Variable                              `gnark:",public"`
	Participant                  frontend.Variable                              `gnark:",public"`
	IncomingPlaceCount           frontend.Variable                              `gnark:",public"`
	IncomingPlaces               [workflow.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaceCount           frontend.Variable                              `gnark:",public"`
	OutgoingPlaces               [workflow.MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount       frontend.Variable `gnark:",public"`
	StartPlace       frontend.Variable `gnark:",public"`
	EndPlace         frontend.Variable `gnark:",public"`
	Transitions      [workflow.MaxTransitionCount]Transition
	ParticipantCount frontend.Variable `gnark:",public"`
}

func FromSaltedHash(commitment crypto.SaltedHash) SaltedHash {
	return SaltedHash{
		Value: commitment.Value,
		Salt:  commitment.Salt,
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

	publicKeyCount := len(instance.PublicKeys)
	if publicKeyCount > workflow.MaxParticipantCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many publicKeys", instance.Id)
	}
	var publicKeys [workflow.MaxParticipantCount]eddsa.PublicKey
	for i := 0; i < publicKeyCount; i++ {
		publicKeys[i] = fromPublicKey(instance.PublicKeys[i])
	}
	for i := publicKeyCount; i < workflow.MaxParticipantCount; i++ {
		publicKeys[i] = emptyPublicKey()
	}

	return Instance{
		TokenCounts: tokenCounts,
		PublicKeys:  publicKeys,
	}, nil
}

func fromPublicKey(publicKeyBytes []byte) eddsa.PublicKey {
	var publicKey eddsa.PublicKey
	publicKey.Assign(twistededwards.BN254, publicKeyBytes)
	return publicKey
}

func emptyPublicKey() eddsa.PublicKey {
	var publicKey eddsa.PublicKey
	zeros := make([]byte, 32)
	publicKey.A.X = zeros
	publicKey.A.Y = zeros
	return publicKey
}

func FromPetriNet(petriNet workflow.PetriNet) (PetriNet, error) {
	placeCount := petriNet.PlaceCount
	if placeCount > workflow.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has too many places", petriNet.Id)
	}
	if petriNet.ParticipantCount > workflow.MaxParticipantCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has too many participants", petriNet.Id)
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
	isExecutableByAnyParticipant := 0
	if transition.IsExecutableByAnyParticipant {
		isExecutableByAnyParticipant = 1
	}
	return Transition{
		IncomingPlaceCount:           incomingPlaceCount,
		IncomingPlaces:               incomingPlaces,
		OutgoingPlaceCount:           outgoingPlaceCount,
		OutgoingPlaces:               outgoingPlaces,
		Participant:                  transition.Participant,
		IsExecutableByAnyParticipant: isExecutableByAnyParticipant,
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
		IncomingPlaceCount:           0,
		IncomingPlaces:               incomingPlaces,
		OutgoingPlaceCount:           0,
		OutgoingPlaces:               outgoingPlaces,
		IsExecutableByAnyParticipant: 0,
		Participant:                  0,
	}
}
