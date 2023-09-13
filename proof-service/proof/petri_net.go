package proof

const MaxTransitionCount = 100
const MaxBranchingFactor = 3

type Transition struct {
	IncomingPlaces [MaxBranchingFactor]uint8
	OutgoingPlaces [MaxBranchingFactor]uint8
}

type PetriNet struct {
	PlaceCount  uint8
	StartPlace  uint8
	Transitions [MaxTransitionCount]Transition
}
