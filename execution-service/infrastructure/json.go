package infrastructure

import (
	"encoding/json"
	"proof-service/domain"
	"proof-service/utils"
	"strconv"
)

type Transition struct {
	Id             string `json:"id"`
	Name           string `json:"name"`
	IncomingPlaces []uint `json:"incomingPlaces"`
	OutgoingPlaces []uint `json:"outgoingPlaces"`
	Participant    string `json:"participant"`
	Message        string `json:"message"`
}

type Model struct {
	Id               string       `json:"id"`
	PlaceCount       uint         `json:"placeCount"`
	ParticipantCount uint         `json:"participantCount"`
	MessageCount     uint         `json:"messageCount"`
	StartPlace       uint         `json:"startPlace"`
	EndPlaces        []uint       `json:"endPlaces"`
	Transitions      []Transition `json:"transitions"`
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
	domainEndPlaces := make([]uint, len(model.EndPlaces))
	copy(domainEndPlaces, model.EndPlaces)
	return domain.Model{
		Id:               model.Id,
		PlaceCount:       model.PlaceCount,
		ParticipantCount: model.ParticipantCount,
		MessageCount:     model.MessageCount,
		StartPlace:       model.StartPlace,
		EndPlaces:        domainEndPlaces,
		Transitions:      domainTransitions,
	}
}

func toDomainTransition(transition Transition) domain.Transition {
	participant := 0
	particpantIsValid := transition.Participant != ""
	if particpantIsValid {
		var err error
		participant, err = strconv.Atoi(transition.Participant)
		utils.PanicOnError(err)
	}
	message := 0
	messageIsValid := transition.Message != ""
	if messageIsValid {
		var err error
		message, err = strconv.Atoi(transition.Message)
		utils.PanicOnError(err)
	}
	return domain.Transition{
		Id:                 transition.Id,
		Name:               transition.Name,
		IncomingPlaces:     transition.IncomingPlaces,
		OutgoingPlaces:     transition.OutgoingPlaces,
		ParticipantIsValid: particpantIsValid,
		Participant:        uint(participant),
		MessageIsValid:     messageIsValid,
		Message:            uint(message),
	}
}
