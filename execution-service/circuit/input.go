package circuit

import (
	"encoding/hex"
	"fmt"
	"proof-service/authentication"
	"proof-service/domain"

	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type Signature struct {
	Value     eddsa.Signature
	PublicKey eddsa.PublicKey
}

type MessageHash struct {
	Value [domain.MessageHashLength]uints.U8
}

type Instance struct {
	Hash          frontend.Variable `gnark:",public"`
	TokenCounts   [domain.MaxPlaceCount]frontend.Variable
	PublicKeys    [domain.MaxParticipantCount]eddsa.PublicKey
	MessageHashes [domain.MaxMessageCount]MessageHash
	Salt          frontend.Variable
}

type Transition struct {
	IncomingPlaceCount frontend.Variable                            `gnark:",public"`
	IncomingPlaces     [domain.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	OutgoingPlaceCount frontend.Variable                            `gnark:",public"`
	OutgoingPlaces     [domain.MaxBranchingFactor]frontend.Variable `gnark:",public"`
	ParticipantIsValid frontend.Variable                            `gnark:",public"`
	Participant        frontend.Variable                            `gnark:",public"`
	MessageIsValid     frontend.Variable                            `gnark:",public"`
	Message            frontend.Variable                            `gnark:",public"`
}

type Model struct {
	PlaceCount       frontend.Variable                          `gnark:",public"`
	StartPlace       frontend.Variable                          `gnark:",public"`
	EndPlaces        [domain.MaxEndPlaceCount]frontend.Variable `gnark:",public"`
	Transitions      [domain.MaxTransitionCount]Transition
	ParticipantCount frontend.Variable `gnark:",public"`
}

func FromSignature(signature authentication.Signature) Signature {
	var value eddsa.Signature
	value.Assign(twistededwards.BN254, signature.Value)
	var publicKey eddsa.PublicKey
	publicKey.Assign(twistededwards.BN254, signature.PublicKey.Value)
	return Signature{
		Value:     value,
		PublicKey: publicKey,
	}
}

func FromInstance(instance domain.Instance) (Instance, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > domain.MaxPlaceCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many places", hex.EncodeToString(instance.Hash))
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
		return Instance{}, fmt.Errorf("instance '%s' has too many publicKeys", hex.EncodeToString(instance.Hash))
	}
	var publicKeys [domain.MaxParticipantCount]eddsa.PublicKey
	for i := 0; i < publicKeyCount; i++ {
		publicKeys[i] = fromPublicKey(instance.PublicKeys[i])
	}
	for i := publicKeyCount; i < domain.MaxParticipantCount; i++ {
		publicKeys[i] = emptyPublicKey()
	}

	messageHashCount := len(instance.MessageHashes)
	if messageHashCount > domain.MaxMessageCount {
		return Instance{}, fmt.Errorf("instance '%s' has too many messageHashes", hex.EncodeToString(instance.Hash))
	}
	var messageHashes [domain.MaxMessageCount]MessageHash
	for i := 0; i < messageHashCount; i++ {
		messageHashes[i] = fromMessageHash(instance.MessageHashes[i])
	}
	for i := messageHashCount; i < domain.MaxMessageCount; i++ {
		messageHashes[i] = emptyMessageHash()
	}

	return Instance{
		Hash:          instance.Hash,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
		Salt:          instance.Salt,
	}, nil
}

func fromPublicKey(publicKey domain.PublicKey) eddsa.PublicKey {
	var eddsaPublicKey eddsa.PublicKey
	eddsaPublicKey.Assign(twistededwards.BN254, publicKey.Value)
	return eddsaPublicKey
}

func emptyPublicKey() eddsa.PublicKey {
	var publicKey eddsa.PublicKey
	zeros := make([]byte, 32)
	publicKey.A.X = zeros
	publicKey.A.Y = zeros
	return publicKey
}

func fromMessageHash(messageHash domain.MessageHash) MessageHash {
	return MessageHash{
		Value: ([domain.MessageHashLength]uints.U8)(uints.NewU8Array(messageHash.Value[:])),
	}
}

func emptyMessageHash() MessageHash {
	var zeros [domain.MessageHashLength]byte
	return MessageHash{
		Value: ([domain.MessageHashLength]uints.U8)(uints.NewU8Array(zeros[:])),
	}
}

func FromModel(model domain.Model) (Model, error) {
	placeCount := model.PlaceCount
	if placeCount > domain.MaxPlaceCount {
		return Model{}, fmt.Errorf("model '%s' has too many places", model.Id)
	}
	if model.ParticipantCount > domain.MaxParticipantCount {
		return Model{}, fmt.Errorf("model '%s' has too many participants", model.Id)
	}
	if model.MessageCount > domain.MaxMessageCount {
		return Model{}, fmt.Errorf("model '%s' has too many messages", model.Id)
	}
	transitions, err := fromTransitions(model.Id, model.Transitions)
	if err != nil {
		return Model{}, err
	}
	if model.StartPlace >= domain.MaxPlaceCount {
		return Model{}, fmt.Errorf("model '%s' has invalid startPlace", model.Id)
	}
	if len(model.EndPlaces) > domain.MaxEndPlaceCount || len(model.EndPlaces) < 1 {
		return Model{}, fmt.Errorf("model '%s' has invalid number of endPlaces", model.Id)
	}
	var endPlaces [domain.MaxEndPlaceCount]frontend.Variable
	for i, endPlace := range model.EndPlaces {
		if endPlace >= domain.MaxPlaceCount {
			return Model{}, fmt.Errorf("model '%s' has invalid endPlace", model.Id)
		}
		endPlaces[i] = endPlace
	}
	for i := len(model.EndPlaces); i < domain.MaxEndPlaceCount; i++ {
		endPlaces[i] = endPlaces[0]
	}
	return Model{
		PlaceCount:       model.PlaceCount,
		StartPlace:       model.StartPlace,
		EndPlaces:        endPlaces,
		Transitions:      transitions,
		ParticipantCount: model.ParticipantCount,
	}, nil
}

func fromTransitions(modelId string, workflowTransitions []domain.Transition) ([domain.MaxTransitionCount]Transition, error) {
	transitionCount := len(workflowTransitions)
	if transitionCount > domain.MaxTransitionCount {
		return [domain.MaxTransitionCount]Transition{}, fmt.Errorf("model '%s' has too many transitions", modelId)
	}
	var transitions [domain.MaxTransitionCount]Transition
	var err error
	for i := 0; i < transitionCount; i++ {
		transitions[i], err = fromTransition(workflowTransitions[i])
		if err != nil {
			return [domain.MaxTransitionCount]Transition{}, fmt.Errorf("model '%s' cannot be mapped because transition at index %d is invalid: %w", modelId, i, err)
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
	return Transition{
		IncomingPlaceCount: incomingPlaceCount,
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: outgoingPlaceCount,
		OutgoingPlaces:     outgoingPlaces,
		Participant:        transition.Participant,
		ParticipantIsValid: boolToInt(transition.ParticipantIsValid),
		Message:            transition.Message,
		MessageIsValid:     boolToInt(transition.MessageIsValid),
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
		IncomingPlaceCount: 0,
		IncomingPlaces:     incomingPlaces,
		OutgoingPlaceCount: 0,
		OutgoingPlaces:     outgoingPlaces,
		ParticipantIsValid: 0,
		Participant:        0,
		MessageIsValid:     0,
		Message:            0,
	}
}

func boolToInt(value bool) int {
	result := 0
	if value {
		result = 1
	}
	return result
}
