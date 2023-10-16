package model

import (
	"execution-service/domain"
	"execution-service/hash"
	"fmt"
	"time"
)

type ConstraintJson struct {
	Coefficients       []int  `json:"coefficients"`
	MessageIds         []uint `json:"messageIds"`
	Offset             int    `json:"offset"`
	ComparisonOperator uint   `json:"comparisonOperator"`
}

type TransitionJson struct {
	Id             string          `json:"id"`
	Name           string          `json:"name"`
	IncomingPlaces []uint          `json:"incomingPlaces"`
	OutgoingPlaces []uint          `json:"outgoingPlaces"`
	Sender         *uint           `json:"sender,omitempty"`
	Recipient      *uint           `json:"recipient,omitempty"`
	Message        *uint           `json:"message,omitempty"`
	Contraint      *ConstraintJson `json:"constraint,omitempty"`
}

type ModelJson struct {
	Hash             hash.SaltedHashJson `json:"hash"`
	Source           string              `json:"source"`
	PlaceCount       uint                `json:"placeCount"`
	ParticipantCount uint                `json:"participantCount"`
	MessageCount     uint                `json:"messageCount"`
	StartPlaces      []uint              `json:"startPlaces"`
	EndPlaces        []uint              `json:"endPlaces"`
	Transitions      []TransitionJson    `json:"transitions"`
	CreatedAt        time.Time           `json:"createdAt"`
}

func (model *ModelJson) id() string {
	return model.Hash.Hash
}

func (model *ModelJson) ToModel() (domain.Model, error) {
	if model.PlaceCount > domain.MaxPlaceCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many places", model.id())
	}
	if model.ParticipantCount > domain.MaxParticipantCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many participants", model.id())
	}
	if model.MessageCount > domain.MaxMessageCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many messages", model.id())
	}
	startPlaceCount := len(model.StartPlaces)
	if startPlaceCount > domain.MaxStartPlaceCount || startPlaceCount < 1 {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid number of startPlaces", model.id())
	}
	startPlaces := make([]domain.PlaceId, startPlaceCount)
	for i, startPlace := range model.StartPlaces {
		if startPlace >= domain.MaxPlaceCount {
			return domain.Model{}, fmt.Errorf("model '%s' has invalid startPlace", model.id())
		}
		startPlaces[i] = domain.PlaceId(startPlace)
	}
	endPlaceCount := len(model.EndPlaces)
	if endPlaceCount > domain.MaxEndPlaceCount || endPlaceCount < 1 {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid number of endPlaces", model.id())
	}
	endPlaces := make([]domain.PlaceId, endPlaceCount)
	for i, endPlace := range model.EndPlaces {
		if endPlace >= domain.MaxPlaceCount {
			return domain.Model{}, fmt.Errorf("model '%s' has invalid endPlace", model.id())
		}
		endPlaces[i] = domain.PlaceId(endPlace)
	}
	transitionCount := len(model.Transitions)
	if transitionCount > domain.MaxTransitionCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many transitions", model.id())
	}
	transitions := make([]domain.Transition, transitionCount)
	for i, transition := range model.Transitions {
		var err error
		transitions[i], err = transition.toTransition()
		if err != nil {
			return domain.Model{}, err
		}
	}
	hash, err := model.Hash.ToHash()
	if err != nil {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid hash", model.id())
	}
	return domain.Model{
		Hash:             hash,
		Source:           model.Source,
		PlaceCount:       uint16(model.PlaceCount),
		ParticipantCount: uint16(model.ParticipantCount),
		MessageCount:     uint16(model.MessageCount),
		StartPlaces:      startPlaces,
		EndPlaces:        endPlaces,
		Transitions:      transitions,
		CreatedAt:        model.CreatedAt.Unix(),
	}, nil
}

func (transition *TransitionJson) toTransition() (domain.Transition, error) {
	incomingPlaceCount := len(transition.IncomingPlaces)
	outgoingPlaceCount := len(transition.OutgoingPlaces)
	if incomingPlaceCount > domain.MaxBranchingFactor || outgoingPlaceCount > domain.MaxBranchingFactor {
		return domain.Transition{}, fmt.Errorf("transition '%s' branches too much", transition.Id)
	}
	incomingPlaces := make([]domain.PlaceId, incomingPlaceCount)
	outgoingPlaces := make([]domain.PlaceId, outgoingPlaceCount)
	for i := 0; i < incomingPlaceCount; i++ {
		incomingPlaces[i] = domain.PlaceId(transition.IncomingPlaces[i])
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		outgoingPlaces[i] = domain.PlaceId(transition.OutgoingPlaces[i])
	}
	initiatingParticipant := domain.EmptyParticipantId
	if transition.Sender != nil {
		initiatingParticipant = uint16(*transition.Sender)
	}
	respondingParticipant := domain.EmptyParticipantId
	if transition.Recipient != nil {
		respondingParticipant = uint16(*transition.Recipient)
	}
	if initiatingParticipant > domain.MaxParticipantCount || respondingParticipant > domain.MaxParticipantCount {
		return domain.Transition{}, fmt.Errorf("transition %s has invalid participant", transition.Id)
	}
	message := domain.EmptyMessageId
	if transition.Message != nil {
		message = uint16(*transition.Message)
	}
	if message > domain.MaxMessageCount {
		return domain.Transition{}, fmt.Errorf("transition %s has invalid message", transition.Id)
	}
	var constraint domain.Constraint
	if transition.Contraint != nil {
		var err error
		constraint, err = transition.Contraint.toConstraint()
		if err != nil {
			return domain.Transition{}, err
		}
	} else {
		constraint = domain.EmptyConstraint()
	}

	return domain.Transition{
		Id:             transition.Id,
		Name:           transition.Name,
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
		Sender:         domain.ParticipantId(initiatingParticipant),
		Recipient:      domain.ParticipantId(respondingParticipant),
		Message:        domain.ModelMessageId(message),
		Constraint:     constraint,
	}, nil
}

func (constraint *ConstraintJson) toConstraint() (domain.Constraint, error) {
	if len(constraint.Coefficients) > domain.MaxConstraintMessageCount {
		return domain.Constraint{}, fmt.Errorf("constraint has too many coefficients")
	}
	if len(constraint.MessageIds) > domain.MaxConstraintMessageCount {
		return domain.Constraint{}, fmt.Errorf("constraint has too many messageIds")
	}
	if len(constraint.MessageIds) != len(constraint.Coefficients) {
		return domain.Constraint{}, fmt.Errorf("number of coefficients differs from number of messageIds in constraint")
	}
	coefficients := make([]domain.IntegerType, len(constraint.Coefficients))
	for i, coefficient := range constraint.Coefficients {
		coefficients[i] = domain.IntegerType(coefficient)
	}
	messageIds := make([]domain.ModelMessageId, len(constraint.MessageIds))
	for i, messageId := range constraint.MessageIds {
		if messageId > domain.MaxMessageCount {
			return domain.Constraint{}, fmt.Errorf("constraint has invalid messageId")
		}
		messageIds[i] = domain.ModelMessageId(messageId)
	}
	if !isValidOparator(constraint.ComparisonOperator) {
		return domain.Constraint{}, fmt.Errorf("constraint has invalid oparator")
	}
	return domain.Constraint{
		Coefficients:       coefficients,
		MessageIds:         messageIds,
		Offset:             domain.IntegerType(constraint.Offset),
		ComparisonOperator: domain.ComparisonOperator(constraint.ComparisonOperator),
	}, nil
}

func ToJson(model domain.Model) ModelJson {
	startPlaces := make([]uint, len(model.StartPlaces))
	for i, startPlace := range model.StartPlaces {
		startPlaces[i] = uint(startPlace)
	}
	endPlaces := make([]uint, len(model.EndPlaces))
	for i, endPlace := range model.EndPlaces {
		endPlaces[i] = uint(endPlace)
	}
	transitions := make([]TransitionJson, len(model.Transitions))
	for i, transition := range model.Transitions {
		transitions[i] = transitionToJson(transition)
	}
	return ModelJson{
		Hash:             hash.ToJson(model.Hash),
		Source:           model.Source,
		PlaceCount:       uint(model.PlaceCount),
		ParticipantCount: uint(model.ParticipantCount),
		MessageCount:     uint(model.MessageCount),
		StartPlaces:      startPlaces,
		EndPlaces:        endPlaces,
		Transitions:      transitions,
		CreatedAt:        time.Unix(model.CreatedAt, 0),
	}
}

func transitionToJson(transition domain.Transition) TransitionJson {
	incomingPlaces := make([]uint, len(transition.IncomingPlaces))
	for i, incomingPlace := range transition.IncomingPlaces {
		incomingPlaces[i] = uint(incomingPlace)
	}
	outgoingPlaces := make([]uint, len(transition.OutgoingPlaces))
	for i, outgoingPlace := range transition.OutgoingPlaces {
		outgoingPlaces[i] = uint(outgoingPlace)
	}
	jsonTransition := TransitionJson{
		Id:             transition.Id,
		Name:           transition.Name,
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
	}
	if transition.Sender != domain.EmptyParticipantId {
		tmp := uint(transition.Sender)
		jsonTransition.Sender = &tmp
	}
	if transition.Recipient != domain.EmptyParticipantId {
		tmp := uint(transition.Recipient)
		jsonTransition.Recipient = &tmp
	}
	if transition.Message != domain.EmptyMessageId {
		tmp := uint(transition.Message)
		jsonTransition.Message = &tmp
	}
	if transition.Constraint.Coefficients != nil {
		tmp := constraintToJson(transition.Constraint)
		jsonTransition.Contraint = &tmp
	}
	return jsonTransition
}

func constraintToJson(constraint domain.Constraint) ConstraintJson {
	coefficients := make([]int, len(constraint.Coefficients))
	for i, coefficient := range constraint.Coefficients {
		coefficients[i] = int(coefficient)
	}
	messageIds := make([]uint, len(constraint.MessageIds))
	for i, messageId := range constraint.MessageIds {
		messageIds[i] = uint(messageId)
	}
	return ConstraintJson{
		Coefficients:       coefficients,
		MessageIds:         messageIds,
		Offset:             int(constraint.Offset),
		ComparisonOperator: uint(constraint.ComparisonOperator),
	}
}

func isValidOparator(oparator uint) bool {
	for _, comparisonOparator := range domain.ValidComparisonOperators {
		if oparator == uint(comparisonOparator) {
			return true
		}
	}
	return false
}
