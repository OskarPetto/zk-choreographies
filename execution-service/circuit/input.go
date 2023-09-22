package circuit

import (
	"encoding/hex"
	"fmt"
	"proof-service/authentication"
	"proof-service/instance"
	"proof-service/model"

	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type Signature struct {
	Value     eddsa.Signature
	PublicKey eddsa.PublicKey
}

type Instance struct {
	Hash        frontend.Variable `gnark:",public"`
	TokenCounts [model.MaxPlaceCount]frontend.Variable
	PublicKeys  [model.MaxParticipantCount]eddsa.PublicKey
	Salt        frontend.Variable
}

type Transition struct {
	IsExecutableByAnyParticipant frontend.Variable                           `gnark:",public"`
	Participant                  frontend.Variable                           `gnark:",public"`
	IncomingPlaceCount           frontend.Variable                           `gnark:",public"`
	IncomingPlaces               [model.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaceCount           frontend.Variable                           `gnark:",public"`
	OutgoingPlaces               [model.MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount       frontend.Variable `gnark:",public"`
	StartPlace       frontend.Variable `gnark:",public"`
	EndPlace         frontend.Variable `gnark:",public"`
	Transitions      [model.MaxTransitionCount]Transition
	ParticipantCount frontend.Variable `gnark:",public"`
}

func FromSignature(signature authentication.Signature) Signature {
	var value eddsa.Signature
	value.Assign(twistededwards.BN254, signature.Value)
	var publicKey eddsa.PublicKey
	publicKey.Assign(twistededwards.BN254, signature.PublicKey)
	return Signature{
		Value:     value,
		PublicKey: publicKey,
	}
}

func FromInstance(instance instance.Instance) (Instance, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > model.MaxPlaceCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many places", hex.EncodeToString(instance.Hash))
	}
	var tokenCounts [model.MaxPlaceCount]frontend.Variable
	for i := 0; i < placeCount; i++ {
		tokenCounts[i] = instance.TokenCounts[i]
	}
	for i := placeCount; i < model.MaxPlaceCount; i++ {
		tokenCounts[i] = 0
	}

	publicKeyCount := len(instance.PublicKeys)
	if publicKeyCount > model.MaxParticipantCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many publicKeys", hex.EncodeToString(instance.Hash))
	}
	var publicKeys [model.MaxParticipantCount]eddsa.PublicKey
	for i := 0; i < publicKeyCount; i++ {
		publicKeys[i] = fromPublicKey(instance.PublicKeys[i])
	}
	for i := publicKeyCount; i < model.MaxParticipantCount; i++ {
		publicKeys[i] = emptyPublicKey()
	}

	return Instance{
		Hash:        instance.Hash,
		TokenCounts: tokenCounts,
		PublicKeys:  publicKeys,
		Salt:        instance.Salt,
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

func FromPetriNet(petriNet model.PetriNet) (PetriNet, error) {
	placeCount := petriNet.PlaceCount
	if placeCount > model.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has too many places", petriNet.Id)
	}
	if petriNet.ParticipantCount > model.MaxParticipantCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has too many participants", petriNet.Id)
	}
	transitions, err := fromTransitions(petriNet.Id, petriNet.Transitions)
	if err != nil {
		return PetriNet{}, err
	}
	if petriNet.StartPlace >= model.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has invalid startPlace", petriNet.Id)
	}
	if petriNet.EndPlace >= model.MaxPlaceCount {
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

func fromTransitions(petriNetId string, workflowTransitions []model.Transition) ([model.MaxTransitionCount]Transition, error) {
	transitionCount := len(workflowTransitions)
	if transitionCount > model.MaxTransitionCount {
		return [model.MaxTransitionCount]Transition{}, fmt.Errorf("petriNet '%s' has too many transitions", petriNetId)
	}
	var transitions [model.MaxTransitionCount]Transition
	var err error
	for i := 0; i < transitionCount; i++ {
		transitions[i], err = fromTransition(workflowTransitions[i])
		if err != nil {
			return [model.MaxTransitionCount]Transition{}, fmt.Errorf("petriNet '%s' cannot be mapped because transition at index %d is invalid: %w", petriNetId, i, err)
		}
	}
	for i := transitionCount; i < model.MaxTransitionCount; i++ {
		transitions[i] = emptyTransition()
	}
	return transitions, nil
}

func fromTransition(transition model.Transition) (Transition, error) {
	incomingPlaceCount := len(transition.IncomingPlaces)
	outgoingPlaceCount := len(transition.OutgoingPlaces)
	if incomingPlaceCount > model.MaxBranchingFactor || outgoingPlaceCount > model.MaxBranchingFactor {
		return Transition{}, fmt.Errorf("transition '%s' branches too much", transition.Id)
	}
	var incomingPlaces [model.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [model.MaxBranchingFactor]frontend.Variable
	for i := 0; i < incomingPlaceCount; i++ {
		incomingPlaces[i] = transition.IncomingPlaces[i]
	}
	for i := incomingPlaceCount; i < model.MaxBranchingFactor; i++ {
		incomingPlaces[i] = model.MaxPlaceCount
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		outgoingPlaces[i] = transition.OutgoingPlaces[i]
	}
	for i := outgoingPlaceCount; i < model.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = model.MaxPlaceCount
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
	var incomingPlaces [model.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [model.MaxBranchingFactor]frontend.Variable
	for i := 0; i < model.MaxBranchingFactor; i++ {
		incomingPlaces[i] = model.MaxPlaceCount
	}
	for i := 0; i < model.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = model.MaxPlaceCount
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
