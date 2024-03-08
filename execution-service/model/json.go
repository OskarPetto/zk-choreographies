package model

import (
	"execution-service/domain"
	"execution-service/hash"
	"execution-service/instance"
	"fmt"
	"time"
)

type ConditionJson struct {
	Coefficients       []int  `json:"coefficients"`
	MessageIds         []uint `json:"messageIds"`
	Offset             int    `json:"offset"`
	ComparisonOperator uint   `json:"comparisonOperator"`
}

type TransitionJson struct {
	Id                    string         `json:"id"`
	Name                  string         `json:"name"`
	IncomingPlaces        []uint         `json:"incomingPlaces"`
	OutgoingPlaces        []uint         `json:"outgoingPlaces"`
	InitiatingParticipant *uint          `json:"initiatingParticipant,omitempty"`
	RespondingParticipant *uint          `json:"respondingParticipant,omitempty"`
	InitiatingMessage     *uint          `json:"initiatingMessage,omitempty"`
	RespondingMessage     *uint          `json:"respondingMessage,omitempty"`
	Condition             *ConditionJson `json:"condition,omitempty"`
}

type ModelJson struct {
	Id               string              `json:"id"`
	SaltedHash       hash.SaltedHashJson `json:"saltedHash"`
	Source           string              `json:"source"`
	PlaceCount       uint                `json:"placeCount"`
	ParticipantCount uint                `json:"participantCount"`
	MessageCount     uint                `json:"messageCount"`
	StartPlaces      []uint              `json:"startPlaces"`
	EndPlaces        []uint              `json:"endPlaces"`
	Transitions      []TransitionJson    `json:"transitions"`
	CreatedAt        time.Time           `json:"createdAt"`
}

type ImportModelCommandJson struct {
	Model    ModelJson             `json:"model"`
	Instance instance.InstanceJson `json:"instance"`
}

func (model *ModelJson) id() string {
	return model.SaltedHash.Hash
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
	hash, err := model.SaltedHash.ToHash()
	if err != nil {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid hash", model.id())
	}
	return domain.Model{
		SaltedHash:       hash,
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
	if transition.InitiatingParticipant != nil {
		initiatingParticipant = uint16(*transition.InitiatingParticipant)
	}
	respondingParticipant := domain.EmptyParticipantId
	if transition.RespondingParticipant != nil {
		respondingParticipant = uint16(*transition.RespondingParticipant)
	}
	if initiatingParticipant > domain.MaxParticipantCount || respondingParticipant > domain.MaxParticipantCount {
		return domain.Transition{}, fmt.Errorf("transition %s has invalid participant", transition.Id)
	}
	initiatingMessage := domain.EmptyMessageId
	if transition.InitiatingMessage != nil {
		initiatingMessage = uint16(*transition.InitiatingMessage)
	}
	if initiatingMessage > domain.MaxMessageCount {
		return domain.Transition{}, fmt.Errorf("transition %s has invalid initiating message", transition.Id)
	}
	respondingMessage := domain.EmptyMessageId
	if transition.RespondingMessage != nil {
		respondingMessage = uint16(*transition.RespondingMessage)
	}
	if respondingMessage > domain.MaxMessageCount {
		return domain.Transition{}, fmt.Errorf("transition %s has invalid responding message", transition.Id)
	}
	var condition domain.Condition
	if transition.Condition != nil {
		var err error
		condition, err = transition.Condition.toCondition()
		if err != nil {
			return domain.Transition{}, err
		}
	} else {
		condition = domain.EmptyCondition()
	}

	return domain.Transition{
		Id:                    transition.Id,
		Name:                  transition.Name,
		IncomingPlaces:        incomingPlaces,
		OutgoingPlaces:        outgoingPlaces,
		InitiatingParticipant: domain.ParticipantId(initiatingParticipant),
		RespondingParticipant: domain.ParticipantId(respondingParticipant),
		InitiatingMessage:     domain.ModelMessageId(initiatingMessage),
		RespondingMessage:     domain.ModelMessageId(respondingMessage),
		Condition:             condition,
	}, nil
}

func (condition *ConditionJson) toCondition() (domain.Condition, error) {
	if len(condition.Coefficients) > domain.MaxMessageCountInConditions {
		return domain.Condition{}, fmt.Errorf("condition has too many coefficients")
	}
	if len(condition.MessageIds) > domain.MaxMessageCountInConditions {
		return domain.Condition{}, fmt.Errorf("condition has too many messageIds")
	}
	if len(condition.MessageIds) != len(condition.Coefficients) {
		return domain.Condition{}, fmt.Errorf("number of coefficients differs from number of messageIds in condition")
	}
	coefficients := make([]domain.IntegerType, len(condition.Coefficients))
	for i, coefficient := range condition.Coefficients {
		coefficients[i] = domain.IntegerType(coefficient)
	}
	messageIds := make([]domain.ModelMessageId, len(condition.MessageIds))
	for i, messageId := range condition.MessageIds {
		if messageId > domain.MaxMessageCount {
			return domain.Condition{}, fmt.Errorf("condition has invalid messageId")
		}
		messageIds[i] = domain.ModelMessageId(messageId)
	}
	if !isValidOparator(condition.ComparisonOperator) {
		return domain.Condition{}, fmt.Errorf("condition has invalid oparator")
	}
	return domain.Condition{
		Coefficients:       coefficients,
		MessageIds:         messageIds,
		Offset:             domain.IntegerType(condition.Offset),
		ComparisonOperator: domain.ComparisonOperator(condition.ComparisonOperator),
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
		Id:               model.Id(),
		SaltedHash:       hash.ToJson(model.SaltedHash),
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
	if transition.InitiatingParticipant != domain.EmptyParticipantId {
		tmp := uint(transition.InitiatingParticipant)
		jsonTransition.InitiatingParticipant = &tmp
	}
	if transition.RespondingParticipant != domain.EmptyParticipantId {
		tmp := uint(transition.RespondingParticipant)
		jsonTransition.RespondingParticipant = &tmp
	}
	if transition.InitiatingMessage != domain.EmptyMessageId {
		tmp := uint(transition.InitiatingMessage)
		jsonTransition.InitiatingMessage = &tmp
	}
	if transition.RespondingMessage != domain.EmptyMessageId {
		tmp := uint(transition.RespondingMessage)
		jsonTransition.RespondingMessage = &tmp
	}
	if len(transition.Condition.Coefficients) > 0 {
		tmp := conditionToJson(transition.Condition)
		jsonTransition.Condition = &tmp
	}
	return jsonTransition
}

func conditionToJson(condition domain.Condition) ConditionJson {
	coefficients := make([]int, len(condition.Coefficients))
	for i, coefficient := range condition.Coefficients {
		coefficients[i] = int(coefficient)
	}
	messageIds := make([]uint, len(condition.MessageIds))
	for i, messageId := range condition.MessageIds {
		messageIds[i] = uint(messageId)
	}
	return ConditionJson{
		Coefficients:       coefficients,
		MessageIds:         messageIds,
		Offset:             int(condition.Offset),
		ComparisonOperator: uint(condition.ComparisonOperator),
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

func (cmd ImportModelCommandJson) ToModelCommand() (ImportModelCommand, error) {
	model, err := cmd.Model.ToModel()
	if err != nil {
		return ImportModelCommand{}, err
	}
	instance, err := cmd.Instance.ToInstance()
	if err != nil {
		return ImportModelCommand{}, err
	}
	return ImportModelCommand{
		Model:    model,
		Instance: instance,
	}, nil
}
