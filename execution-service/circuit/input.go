package circuit

import (
	"proof-service/authentication"
	"proof-service/domain"
	"proof-service/utils"

	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const defaultMessageHash = "18386210742325734038511415457231681258408421947992479991590796204613365952235"

type Signature struct {
	Value     eddsa.Signature
	PublicKey eddsa.PublicKey
}

type Instance struct {
	Hash          frontend.Variable `gnark:",public"`
	TokenCounts   [domain.MaxPlaceCount]frontend.Variable
	PublicKeys    [domain.MaxParticipantCount]eddsa.PublicKey
	MessageHashes [domain.MaxMessageCount]frontend.Variable
	Salt          frontend.Variable
}

type Transition struct {
	IsInitialized  frontend.Variable
	IncomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	OutgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	Participant    frontend.Variable
	Message        frontend.Variable
}

type Model struct {
	Hash             frontend.Variable `gnark:",public"`
	PlaceCount       frontend.Variable
	ParticipantCount frontend.Variable
	MessageCount     frontend.Variable
	StartPlaces      [domain.MaxStartPlaceCount]frontend.Variable
	EndPlaces        [domain.MaxEndPlaceCount]frontend.Variable
	Transitions      [domain.MaxTransitionCount]Transition
	Salt             frontend.Variable
}

func FromSignature(signature authentication.Signature) Signature {
	var value eddsa.Signature
	value.Assign(twistededwards.BN254, signature.Value)
	var publicKey eddsa.PublicKey
	publicKey.Assign(twistededwards.BN254, signature.PublicKey.Value[:])
	return Signature{
		Value:     value,
		PublicKey: publicKey,
	}
}

func FromInstance(instance domain.Instance) (Instance, error) {
	var tokenCounts [domain.MaxPlaceCount]frontend.Variable
	for i, tokenCount := range instance.TokenCounts {
		tokenCounts[i] = tokenCount
	}
	var publicKeys [domain.MaxParticipantCount]eddsa.PublicKey
	for i, publicKey := range instance.PublicKeys {
		publicKeys[i] = fromPublicKey(publicKey)
	}
	var messageHashes [domain.MaxMessageCount]frontend.Variable
	for i, messageHash := range instance.MessageHashes {
		messageHashes[i] = utils.HashToField(messageHash.Value)
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
	eddsaPublicKey.Assign(twistededwards.BN254, publicKey.Value[:])
	return eddsaPublicKey
}

func FromModel(model domain.Model) (Model, error) {
	var startPlaces [domain.MaxStartPlaceCount]frontend.Variable
	for i, startPlace := range model.StartPlaces {
		startPlaces[i] = startPlace
	}
	var endPlaces [domain.MaxEndPlaceCount]frontend.Variable
	for i, endPlace := range model.EndPlaces {
		endPlaces[i] = endPlace
	}
	var transitions [domain.MaxTransitionCount]Transition
	for i, transition := range model.Transitions {
		transitions[i] = fromTransition(transition)
	}
	return Model{
		Hash:             model.Hash,
		PlaceCount:       model.PlaceCount,
		ParticipantCount: model.ParticipantCount,
		MessageCount:     model.MessageCount,
		StartPlaces:      startPlaces,
		EndPlaces:        endPlaces,
		Transitions:      transitions,
		Salt:             model.Salt,
	}, nil
}

func fromTransition(transition domain.Transition) Transition {
	var incomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i, incomingPlace := range transition.IncomingPlaces {
		incomingPlaces[i] = incomingPlace
	}
	var outgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i, outgoingPlace := range transition.OutgoingPlaces {
		outgoingPlaces[i] = outgoingPlace
	}
	isInitialized := 0
	if transition.IsInitialized {
		isInitialized = 1
	}
	return Transition{
		IsInitialized:  isInitialized,
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
		Participant:    transition.Participant,
		Message:        transition.Message,
	}
}
