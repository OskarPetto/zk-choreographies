package workflow

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
	PlaceCount  uint
	Transitions []Transition
}

type Instance struct {
	Id          string
	TokenCounts []int
}

func SerializeInstance(instance Instance) []byte {
	placeCount := len(instance.TokenCounts)
	var bytes = make([]byte, MaxPlaceCount+1)
	bytes[0] = byte(placeCount)
	for i := 0; i < placeCount; i++ {
		bytes[i+1] = byte(instance.TokenCounts[i])
	}
	return bytes
}
