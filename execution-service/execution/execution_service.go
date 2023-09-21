package execution

import (
	"fmt"
	"proof-service/domain"
	"proof-service/proof"

	"github.com/google/uuid"
)

type ExecutionService struct {
	proofService proof.ProofService
}

func NewExecutionService() ExecutionService {
	return ExecutionService{
		proofService: proof.NewProofService(),
	}
}

func InstantiatePetriNet(petriNet domain.PetriNet, publicKeys [][]byte) (domain.Instance, error) {
	if int(petriNet.ParticipantCount) != len(publicKeys) {
		return domain.Instance{}, fmt.Errorf("the number of public keys does not match the number of participants in the petriNet %s", petriNet.Id)
	}
	tokenCounts := make([]int, petriNet.PlaceCount)
	tokenCounts[petriNet.StartPlace] = 1
	return domain.Instance{
		Id:          createInstanceId(),
		TokenCounts: tokenCounts,
		PublicKeys:  publicKeys,
	}, nil
}

func ExecuteTransition()

func createInstanceId() string {
	return uuid.New().String()
}
