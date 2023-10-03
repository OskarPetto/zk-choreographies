package circuit

import (
	"execution-service/domain"
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type Hash struct {
	Value frontend.Variable `gnark:",public"`
	Salt  frontend.Variable
}

type Signature struct {
	Value     eddsa.Signature
	PublicKey eddsa.PublicKey
}

type Instance struct {
	Hash          Hash
	TokenCounts   [domain.MaxPlaceCount]frontend.Variable
	PublicKeys    [domain.MaxParticipantCount]eddsa.PublicKey
	MessageHashes [domain.MaxMessageCount]frontend.Variable
}

type ConstraintInput struct {
	IntegerMessages [domain.MaxConstraintMessageCount]frontend.Variable
	Salts           [domain.MaxConstraintMessageCount]frontend.Variable
}

type Constraint struct {
	Coefficients       [domain.MaxConstraintMessageCount]frontend.Variable
	MessageIds         [domain.MaxConstraintMessageCount]frontend.Variable
	Offset             frontend.Variable
	ComparisonOperator frontend.Variable
}

type Transition struct {
	IsTransition   frontend.Variable
	IncomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	OutgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	Participant    frontend.Variable
	Message        frontend.Variable
	Constraint     Constraint
}

type Model struct {
	Hash             Hash
	PlaceCount       frontend.Variable
	ParticipantCount frontend.Variable
	MessageCount     frontend.Variable
	StartPlaces      [domain.MaxStartPlaceCount]frontend.Variable
	EndPlaces        [domain.MaxEndPlaceCount]frontend.Variable
	Transitions      [domain.MaxTransitionCount]Transition
}

func FromSignature(signature domain.Signature) Signature {
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
		messageHashes[i] = fromBytes(messageHash)
	}
	return Instance{
		Hash:          fromHash(instance.Hash),
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
		Hash:             fromHash(model.Hash),
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
	isTransition := 0
	if transition.IsTransition {
		isTransition = 1
	}
	return Transition{
		IsTransition:   isTransition,
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
		Participant:    transition.Participant,
		Message:        transition.Message,
		Constraint:     fromConstraint(transition.Constraint),
	}
}

func fromConstraint(constraint domain.Constraint) Constraint {
	var coefficients [domain.MaxConstraintMessageCount]frontend.Variable
	for i, coefficient := range constraint.Coefficients {
		coefficients[i] = coefficient
	}
	var messageIds [domain.MaxConstraintMessageCount]frontend.Variable
	for i, messageId := range constraint.MessageIds {
		messageIds[i] = messageId
	}
	return Constraint{
		Coefficients:       coefficients,
		MessageIds:         messageIds,
		Offset:             constraint.Offset,
		ComparisonOperator: constraint.ComparisonOperator,
	}
}

func FromConstraintInput(input domain.ConstraintInput) ConstraintInput {
	var integerMessages [domain.MaxConstraintMessageCount]frontend.Variable
	var salts [domain.MaxConstraintMessageCount]frontend.Variable
	for i, message := range input.IntegerMessages {
		integerMessages[i] = message.IntegerMessage
		salts[i] = fromBytes(message.Hash.Salt)
	}
	return ConstraintInput{
		IntegerMessages: integerMessages,
		Salts:           salts,
	}
}

func fromHash(hash domain.Hash) Hash {
	return Hash{
		Value: fromBytes(hash.Value),
		Salt:  fromBytes(hash.Salt),
	}
}

func fromBytes(data [fr.Bytes]byte) frontend.Variable {
	fieldElement, err := fr.BigEndian.Element(&data)
	utils.PanicOnError(err)
	return fieldElement
}
