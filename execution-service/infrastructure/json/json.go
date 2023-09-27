package json

import (
	"encoding/json"
	"fmt"
	"proof-service/domain"
)

type Transition struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IncomingPlaces []uint `json:"incomingPlaces"`
	OutgoingPlaces []uint `json:"outgoingPlaces"`
	Participant    uint   `json:"participant"`
	Message        uint   `json:"message"`
}

type Model struct {
	Id               string       `json:"id"`
	PlaceCount       uint         `json:"placeCount"`
	ParticipantCount uint         `json:"participantCount"`
	MessageCount     uint         `json:"messageCount"`
	StartPlaces      []uint       `json:"startPlaces"`
	EndPlaces        []uint       `json:"endPlaces"`
	Transitions      []Transition `json:"transitions"`
}

func (transition *Transition) UnmarshalJSON(data []byte) error {
	type Alias Transition
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
		transition.Participant = domain.DefaultParticipantId
	} else {
		transition.Participant = *(tmp.Participant)
	}
	if tmp.Message == nil {
		transition.Message = domain.DefaultMessageId
	} else {
		transition.Message = *(tmp.Message)
	}
	return nil
}

func UnmarshalModel(data []byte) (domain.Model, error) {
	var model Model
	err := json.Unmarshal(data, &model)
	if err != nil {
		return domain.Model{}, err
	}
	return toDomainModel(model)
}

func toDomainModel(model Model) (domain.Model, error) {
	if model.PlaceCount > domain.MaxPlaceCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many places", model.Id)
	}
	if model.ParticipantCount > domain.MaxParticipantCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many participants", model.Id)
	}
	if model.MessageCount > domain.MaxMessageCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many messages", model.Id)
	}
	startPlaceCount := len(model.StartPlaces)
	if startPlaceCount > domain.MaxStartPlaceCount || startPlaceCount < 1 {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid number of startPlaces", model.Id)
	}
	var startPlaces [domain.MaxStartPlaceCount]domain.PlaceId
	for i, startPlace := range model.StartPlaces {
		if startPlace >= domain.MaxPlaceCount {
			return domain.Model{}, fmt.Errorf("model '%s' has invalid startPlace", model.Id)
		}
		startPlaces[i] = domain.PlaceId(startPlace)
	}
	for i := startPlaceCount; i < domain.MaxStartPlaceCount; i++ {
		startPlaces[i] = domain.DefaultPlaceId
	}
	endPlaceCount := len(model.EndPlaces)
	if endPlaceCount > domain.MaxEndPlaceCount || endPlaceCount < 1 {
		return domain.Model{}, fmt.Errorf("model '%s' has invalid number of endPlaces", model.Id)
	}
	var endPlaces [domain.MaxEndPlaceCount]domain.PlaceId
	for i, endPlace := range model.EndPlaces {
		if endPlace >= domain.MaxPlaceCount {
			return domain.Model{}, fmt.Errorf("model '%s' has invalid endPlace", model.Id)
		}
		endPlaces[i] = domain.PlaceId(endPlace)
	}
	for i := endPlaceCount; i < domain.MaxEndPlaceCount; i++ {
		endPlaces[i] = domain.DefaultPlaceId
	}
	transitionCount := len(model.Transitions)
	if transitionCount > domain.MaxTransitionCount {
		return domain.Model{}, fmt.Errorf("model '%s' has too many transitions", model.Id)
	}
	var transitions [domain.MaxTransitionCount]domain.Transition
	for i, transition := range model.Transitions {
		var err error
		transitions[i], err = toDomainTransition(transition)
		if err != nil {
			return domain.Model{}, err
		}
	}
	for i := transitionCount; i < domain.MaxTransitionCount; i++ {
		transitions[i] = domain.DefaultTransition
	}
	domainModel := domain.Model{
		Id:               model.Id,
		PlaceCount:       uint8(model.PlaceCount),
		ParticipantCount: uint8(model.ParticipantCount),
		MessageCount:     uint8(model.MessageCount),
		StartPlaces:      startPlaces,
		EndPlaces:        endPlaces,
		Transitions:      transitions,
	}
	domainModel.ComputeHash()
	return domainModel, nil
}

func toDomainTransition(transition Transition) (domain.Transition, error) {
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
		incomingPlaces[i] = domain.DefaultPlaceId
	}
	for i := 0; i < outgoingPlaceCount; i++ {
		outgoingPlaces[i] = domain.PlaceId(transition.OutgoingPlaces[i])
	}
	for i := outgoingPlaceCount; i < domain.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = domain.DefaultPlaceId
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
		IsInitialized:  true,
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
		Participant:    domain.ParticipantId(transition.Participant),
		Message:        domain.MessageId(transition.Message),
	}, nil
}
