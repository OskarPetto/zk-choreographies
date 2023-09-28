package model

import (
	"encoding/json"
	"execution-service/domain"
	"execution-service/hash"
	"fmt"
	"time"
)

type TransitionJson struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IncomingPlaces []uint `json:"incomingPlaces"`
	OutgoingPlaces []uint `json:"outgoingPlaces"`
	Participant    uint   `json:"participant,omitempty"`
	Message        uint   `json:"message,omitempty"`
}

type ModelJson struct {
	Hash             hash.HashJson    `json:"hash"`
	Choreography     string           `json:"name"`
	PlaceCount       uint             `json:"placeCount"`
	ParticipantCount uint             `json:"participantCount"`
	MessageCount     uint             `json:"messageCount"`
	StartPlaces      []uint           `json:"startPlaces"`
	EndPlaces        []uint           `json:"endPlaces"`
	Transitions      []TransitionJson `json:"transitions"`
	CreatedAt        time.Time        `json:"createdAt"`
}

func (model *ModelJson) Id() string {
	return model.Hash.Value
}

func (transition *TransitionJson) UnmarshalJSON(data []byte) error {
	type Alias TransitionJson
	tmp := struct {
		*Alias
		Participant *uint `json:"participant"`
		Message     *uint `json:"message"`
	}{
		Alias: (*Alias)(transition),
	}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	if tmp.Participant == nil {
		transition.Participant = domain.InvalidParticipantId
	} else {
		transition.Participant = *(tmp.Participant)
	}
	if tmp.Message == nil {
		transition.Message = domain.InvalidMessageId
	} else {
		transition.Message = *(tmp.Message)
	}
	return nil
}

func FromJson(data []byte) (domain.Model, error) {
	var model ModelJson
	err := json.Unmarshal(data, &model)
	if err != nil {
		return domain.Model{}, err
	}
	return model.ToModel()
}

func (model *ModelJson) ToModel() (domain.Model, error) {
	if model.PlaceCount > domain.MaxPlaceCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many places", model.Id())
	}
	if model.ParticipantCount > domain.MaxParticipantCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many participants", model.Id())
	}
	if model.MessageCount > domain.MaxMessageCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many messages", model.Id())
	}
	startPlaceCount := len(model.StartPlaces)
	if startPlaceCount > domain.MaxStartPlaceCount || startPlaceCount < 1 {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid number of startPlaces", model.Id())
	}
	var startPlaces [domain.MaxStartPlaceCount]domain.PlaceId
	for i, startPlace := range model.StartPlaces {
		if startPlace >= domain.MaxPlaceCount {
			return domain.Model{}, fmt.Errorf("model '%s' has invalid startPlace", model.Id())
		}
		startPlaces[i] = domain.PlaceId(startPlace)
	}
	for i := startPlaceCount; i < domain.MaxStartPlaceCount; i++ {
		startPlaces[i] = domain.InvalidPlaceId
	}
	endPlaceCount := len(model.EndPlaces)
	if endPlaceCount > domain.MaxEndPlaceCount || endPlaceCount < 1 {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid number of endPlaces", model.Id())
	}
	var endPlaces [domain.MaxEndPlaceCount]domain.PlaceId
	for i, endPlace := range model.EndPlaces {
		if endPlace >= domain.MaxPlaceCount {
			return domain.Model{}, fmt.Errorf("model '%s' has invalid endPlace", model.Id())
		}
		endPlaces[i] = domain.PlaceId(endPlace)
	}
	for i := endPlaceCount; i < domain.MaxEndPlaceCount; i++ {
		endPlaces[i] = domain.InvalidPlaceId
	}
	transitionCount := len(model.Transitions)
	if transitionCount > domain.MaxTransitionCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many transitions", model.Id())
	}
	var transitions [domain.MaxTransitionCount]domain.Transition
	for i, transition := range model.Transitions {
		var err error
		transitions[i], err = transition.toTransition()
		if err != nil {
			return domain.Model{}, err
		}
	}
	for i := transitionCount; i < domain.MaxTransitionCount; i++ {
		transitions[i] = domain.InvalidTransition()
	}
	hash, err := model.Hash.ToHash()
	if err != nil {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid hash", model.Id())
	}
	return domain.Model{
		Hash:             hash,
		Choreography:     model.Choreography,
		PlaceCount:       uint8(model.PlaceCount),
		ParticipantCount: uint8(model.ParticipantCount),
		MessageCount:     uint8(model.MessageCount),
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
	var incomingPlaces [domain.MaxBranchingFactor]domain.PlaceId
	var outgoingPlaces [domain.MaxBranchingFactor]domain.PlaceId
	for i := 0; i < incomingPlaceCount; i++ {
		incomingPlaces[i] = domain.PlaceId(transition.IncomingPlaces[i])
	}
	for i := incomingPlaceCount; i < domain.MaxBranchingFactor; i++ {
		incomingPlaces[i] = domain.InvalidPlaceId
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		outgoingPlaces[i] = domain.PlaceId(transition.OutgoingPlaces[i])
	}
	for i := outgoingPlaceCount; i < domain.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = domain.InvalidPlaceId
	}
	if transition.Participant > domain.MaxParticipantCount {
		return domain.Transition{}, fmt.Errorf("transition %s has invalid participant", transition.Id)
	}
	if transition.Message > domain.MaxMessageCount {
		return domain.Transition{}, fmt.Errorf("transition %s has invalid message", transition.Id)
	}

	return domain.Transition{
		Id:             transition.Id,
		Name:           transition.Name,
		IsValid:        true,
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
		Participant:    domain.ParticipantId(transition.Participant),
		Message:        domain.MessageId(transition.Message),
	}, nil
}

func ToJson(model domain.Model) ModelJson {
	startPlaces := make([]uint, 0)
	for _, startPlace := range model.StartPlaces {
		if startPlace != domain.InvalidPlaceId {
			break
		}
		startPlaces = append(startPlaces, uint(startPlace))
	}
	endPlaces := make([]uint, 0)
	for _, endPlace := range model.EndPlaces {
		if endPlace == domain.InvalidPlaceId {
			break
		}
		endPlaces = append(endPlaces, uint(endPlace))
	}
	transitions := make([]TransitionJson, 0)
	for _, transition := range model.Transitions {
		if !transition.IsValid {
			break
		}
		transitions = append(transitions, transitionToJson(transition))
	}
	return ModelJson{
		Hash:             hash.HashToJson(model.Hash),
		Choreography:     model.Choreography,
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
	incomingPlaces := make([]uint, 0)
	for _, incomingPlace := range transition.IncomingPlaces {
		if incomingPlace == domain.InvalidPlaceId {
			break
		}
		incomingPlaces = append(incomingPlaces, uint(incomingPlace))
	}
	outgoingPlaces := make([]uint, 0)
	for _, outgoingPlace := range transition.OutgoingPlaces {
		if outgoingPlace == domain.InvalidPlaceId {
			break
		}
		outgoingPlaces = append(outgoingPlaces, uint(outgoingPlace))
	}
	jsonTransition := TransitionJson{
		Id:             transition.Id,
		Name:           transition.Name,
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
	}
	if transition.Participant != domain.InvalidParticipantId {
		jsonTransition.Participant = uint(transition.Participant)
	}
	if transition.Message != domain.InvalidMessageId {
		jsonTransition.Message = uint(transition.Message)
	}
	return jsonTransition
}
