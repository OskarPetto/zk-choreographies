package instance

import (
	"fmt"
	"proof-service/model"
)

type Instance struct {
	Hash        []byte
	TokenCounts []int
	PublicKeys  [][]byte
	Salt        []byte
}

func InstantiatePetriNet(petriNet model.PetriNet, publicKeys [][]byte) (Instance, error) {
	if int(petriNet.ParticipantCount) != len(publicKeys) {
		return Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the petriNet %s", petriNet.Id)
	}
	tokenCounts := make([]int, petriNet.PlaceCount)
	tokenCounts[petriNet.StartPlace] = 1
	instance := Instance{
		TokenCounts: tokenCounts,
		PublicKeys:  publicKeys,
	}
	instance.ComputeHash()
	return instance, nil
}

func (instance Instance) ExecuteTransition(transition model.Transition) (Instance, error) {
	if !isTransitionExecutable(instance, transition) {
		return instance, fmt.Errorf("transition %s is not executable", transition.Id)
	}
	tokenCounts := make([]int, len(instance.TokenCounts))
	copy(tokenCounts, instance.TokenCounts)
	for _, incomingPlaceId := range transition.IncomingPlaces {
		tokenCounts[incomingPlaceId] -= 1
	}
	for _, outgoingPlaceId := range transition.OutgoingPlaces {
		tokenCounts[outgoingPlaceId] += 1
	}
	instance.TokenCounts = tokenCounts
	instance.ComputeHash()
	return instance, nil
}

func isTransitionExecutable(instance Instance, transition model.Transition) bool {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if instance.TokenCounts[incomingPlaceId] < 1 {
			return false
		}
	}
	return true
}
