package circuit

import (
	"fmt"
	"proof-service/crypto"
	"proof-service/domain"

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
	TokenCounts [domain.MaxPlaceCount]frontend.Variable
	PublicKeys  [domain.MaxParticipantCount]eddsa.PublicKey
}

type Transition struct {
	IsExecutableByAnyParticipant frontend.Variable                            `gnark:",public"`
	Participant                  frontend.Variable                            `gnark:",public"`
	IncomingPlaceCount           frontend.Variable                            `gnark:",public"`
	IncomingPlaces               [domain.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaceCount           frontend.Variable                            `gnark:",public"`
	OutgoingPlaces               [domain.MaxBranchingFactor]frontend.Variable `gnark:",public"`
}

type PetriNet struct {
	PlaceCount       frontend.Variable `gnark:",public"`
	StartPlace       frontend.Variable `gnark:",public"`
	EndPlace         frontend.Variable `gnark:",public"`
	Transitions      [domain.MaxTransitionCount]Transition
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

func FromInstance(instance domain.Instance) (Instance, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > domain.MaxPlaceCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many places", instance.Id)
	}
	var tokenCounts [domain.MaxPlaceCount]frontend.Variable
	for i := 0; i < placeCount; i++ {
		tokenCounts[i] = instance.TokenCounts[i]
	}
	for i := placeCount; i < domain.MaxPlaceCount; i++ {
		tokenCounts[i] = 0
	}

	publicKeyCount := len(instance.PublicKeys)
	if publicKeyCount > domain.MaxParticipantCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many publicKeys", instance.Id)
	}
	var publicKeys [domain.MaxParticipantCount]eddsa.PublicKey
	for i := 0; i < publicKeyCount; i++ {
		publicKeys[i] = fromPublicKey(instance.PublicKeys[i])
	}
	for i := publicKeyCount; i < domain.MaxParticipantCount; i++ {
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

func FromPetriNet(petriNet domain.PetriNet) (PetriNet, error) {
	placeCount := petriNet.PlaceCount
	if placeCount > domain.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has too many places", petriNet.Id)
	}
	if petriNet.ParticipantCount > domain.MaxParticipantCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has too many participants", petriNet.Id)
	}
	transitions, err := fromTransitions(petriNet.Id, petriNet.Transitions)
	if err != nil {
		return PetriNet{}, err
	}
	if petriNet.StartPlace >= domain.MaxPlaceCount {
		return PetriNet{}, fmt.Errorf("petriNet '%s' has invalid startPlace", petriNet.Id)
	}
	if petriNet.EndPlace >= domain.MaxPlaceCount {
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

func fromTransitions(petriNetId string, workflowTransitions []domain.Transition) ([domain.MaxTransitionCount]Transition, error) {
	transitionCount := len(workflowTransitions)
	if transitionCount > domain.MaxTransitionCount {
		return [domain.MaxTransitionCount]Transition{}, fmt.Errorf("petriNet '%s' has too many transitions", petriNetId)
	}
	var transitions [domain.MaxTransitionCount]Transition
	var err error
	for i := 0; i < transitionCount; i++ {
		transitions[i], err = fromTransition(workflowTransitions[i])
		if err != nil {
			return [domain.MaxTransitionCount]Transition{}, fmt.Errorf("petriNet '%s' cannot be mapped because transition at index %d is invalid: %w", petriNetId, i, err)
		}
	}
	for i := transitionCount; i < domain.MaxTransitionCount; i++ {
		transitions[i] = emptyTransition()
	}
	return transitions, nil
}

func fromTransition(transition domain.Transition) (Transition, error) {
	incomingPlaceCount := len(transition.IncomingPlaces)
	outgoingPlaceCount := len(transition.OutgoingPlaces)
	if incomingPlaceCount > domain.MaxBranchingFactor || outgoingPlaceCount > domain.MaxBranchingFactor {
		return Transition{}, fmt.Errorf("transition '%s' branches too much", transition.Id)
	}
	var incomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i := 0; i < incomingPlaceCount; i++ {
		incomingPlaces[i] = transition.IncomingPlaces[i]
	}
	for i := incomingPlaceCount; i < domain.MaxBranchingFactor; i++ {
		incomingPlaces[i] = domain.MaxPlaceCount
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		outgoingPlaces[i] = transition.OutgoingPlaces[i]
	}
	for i := outgoingPlaceCount; i < domain.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = domain.MaxPlaceCount
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
	var incomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	var outgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i := 0; i < domain.MaxBranchingFactor; i++ {
		incomingPlaces[i] = domain.MaxPlaceCount
	}
	for i := 0; i < domain.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = domain.MaxPlaceCount
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
