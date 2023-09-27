package circuit

import (
	"execution-service/authentication"
	"execution-service/domain"
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

const defaultMessageHash = 0

type Signature struct {
	Value     eddsa.Signature
	PublicKey eddsa.PublicKey
}

type Instance struct {
	Hash          frontend.Variable `gnark:",public"`
	Salt          frontend.Variable
	TokenCounts   [domain.MaxPlaceCount]frontend.Variable
	PublicKeys    [domain.MaxParticipantCount]eddsa.PublicKey
	MessageHashes [domain.MaxMessageCount]frontend.Variable
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
	Salt             frontend.Variable
	PlaceCount       frontend.Variable
	ParticipantCount frontend.Variable
	MessageCount     frontend.Variable
	StartPlaces      [domain.MaxStartPlaceCount]frontend.Variable
	EndPlaces        [domain.MaxEndPlaceCount]frontend.Variable
	Transitions      [domain.MaxTransitionCount]Transition
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

func FromInstance(instance domain.Instance) Instance {
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
		messageHashes[i] = fromBytes(messageHash.Value)
	}
	return Instance{
		Hash:          fromBytes(instance.Hash.Value),
		Salt:          fromBytes(instance.Hash.Salt),
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
	}
}

func fromPublicKey(publicKey domain.PublicKey) eddsa.PublicKey {
	var eddsaPublicKey eddsa.PublicKey
	eddsaPublicKey.Assign(twistededwards.BN254, publicKey.Value[:])
	return eddsaPublicKey
}

func FromModel(model domain.Model) Model {
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
		Hash:             fromBytes(model.Hash.Value),
		Salt:             fromBytes(model.Hash.Salt),
		PlaceCount:       model.PlaceCount,
		ParticipantCount: model.ParticipantCount,
		MessageCount:     model.MessageCount,
		StartPlaces:      startPlaces,
		EndPlaces:        endPlaces,
		Transitions:      transitions,
	}
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

func fromBytes(data [fr.Bytes]byte) frontend.Variable {
	fieldElement, err := fr.BigEndian.Element(&data)
	utils.PanicOnError(err)
	return fieldElement
}
