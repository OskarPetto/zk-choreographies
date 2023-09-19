package workflow

import "fmt"

const MaxPlaceCount = 100
const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type PlaceId = uint

type Transition struct {
	Id             string
	IncomingPlaces []PlaceId
	OutgoingPlaces []PlaceId
}

type PetriNet struct {
	Id          string
	StartPlace  PlaceId
	EndPlace    PlaceId
	PlaceCount  uint
	Transitions []Transition
}

type Instance struct {
	Id          string
	TokenCounts []int
}

func SerializeInstance(instance Instance) ([]byte, error) {
	placeCount := len(instance.TokenCounts)
	if placeCount > MaxPlaceCount {
		return []byte{}, fmt.Errorf("instance '%s' is too large", instance.Id)
	}
	var bytes = make([]byte, MaxPlaceCount+1)
	bytes[0] = byte(placeCount)
	for i := 0; i < placeCount; i++ {
		bytes[i+1] = byte(instance.TokenCounts[i])
	}
	return bytes, nil
}
