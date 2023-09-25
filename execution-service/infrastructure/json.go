package infrastructure

import (
	"encoding/json"
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
		transition.Participant = domain.MaxParticipantCount
	} else {
		transition.Participant = *(tmp.Participant)
	}
	if tmp.Message == nil {
		transition.Message = domain.MaxMessageCount
	} else {
		transition.Message = *(tmp.Message)
	}
	return nil
}

func FromJson(data []byte) domain.Model {
	var model Model
	json.Unmarshal(data, &model)
	return toDomainModel(model)
}

func toDomainModel(model Model) domain.Model {
	domainTransitions := make([]domain.Transition, len(model.Transitions))
	for i, transition := range model.Transitions {
		domainTransitions[i] = toDomainTransition(transition)
	}
	return domain.Model{
		Id:               model.Id,
		PlaceCount:       model.PlaceCount,
		ParticipantCount: model.ParticipantCount,
		MessageCount:     model.MessageCount,
		StartPlaces:      model.StartPlaces,
		EndPlaces:        model.EndPlaces,
		Transitions:      domainTransitions,
	}
}

func toDomainTransition(transition Transition) domain.Transition {
	return domain.Transition{
		Id:             transition.Id,
		Name:           transition.Name,
		IncomingPlaces: transition.IncomingPlaces,
		OutgoingPlaces: transition.OutgoingPlaces,
		Participant:    transition.Participant,
		Message:        transition.Message,
	}
}
