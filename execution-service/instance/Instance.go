package instance

type InstanceId string

type Instance struct {
	Id          InstanceId
	TokenCounts []int8
}
