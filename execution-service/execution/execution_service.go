package execution

import (
	"fmt"
	"proof-service/domain"
	"proof-service/proof"

	"github.com/google/uuid"
)

type ExecutionService struct {
	isLoaded     bool
	proofService proof.ProofService
}

var executionService ExecutionService

func NewExecutionService() ExecutionService {
	if !executionService.isLoaded {
		executionService = ExecutionService{
			isLoaded:     true,
			proofService: proof.NewProofService(),
		}
	}
	return executionService
}

func (service *ExecutionService) InstantiatePetriNet(petriNet domain.PetriNet, publicKeys [][]byte) (domain.Instance, error) {
	if int(petriNet.ParticipantCount) != len(publicKeys) {
		return domain.Instance{}, fmt.Errorf("the number of public keys must match the number of participants in the petriNet %s", petriNet.Id)
	}
	tokenCounts := make([]int, petriNet.PlaceCount)
	tokenCounts[petriNet.StartPlace] = 1
	return domain.Instance{
		Id:          createInstanceId(),
		TokenCounts: tokenCounts,
		PublicKeys:  publicKeys,
	}, nil
}

func (service *ExecutionService) ExecuteTransition(instance domain.Instance, transition domain.Transition) (domain.Instance, error) {
	newInstance := instance
	if !isTransitionExecutable(instance, transition) {
		return newInstance, fmt.Errorf("transition %s is not executable", transition.Id)
	}
	for _, incomingPlaceId := range transition.IncomingPlaces {
		newInstance.TokenCounts[incomingPlaceId] -= 1
	}
	for _, outgoingPlaceId := range transition.OutgoingPlaces {
		newInstance.TokenCounts[outgoingPlaceId] += 1
	}
	return newInstance, nil
}

func isTransitionExecutable(instance domain.Instance, transition domain.Transition) bool {
	for _, incomingPlaceId := range transition.IncomingPlaces {
		if instance.TokenCounts[incomingPlaceId] < 1 {
			return false
		}
	}
	return true
}

func createInstanceId() string {
	return uuid.New().String()
}
