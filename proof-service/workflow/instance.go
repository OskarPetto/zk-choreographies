package workflow

type InstanceId string

type Instance struct {
	Id          InstanceId
	PetriNet    PetriNetId
	TokenCounts []int8
}
