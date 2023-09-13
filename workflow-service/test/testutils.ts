import {
  PlaceId,
  PetriNet,
  TransitionId,
  Transition,
} from 'src/model/petri-net/petri-net';

export function logObject(object: any) {
  console.dir(object, { depth: null });
}

export function findTransitions(
  petriNet: PetriNet,
  transitionIds: TransitionId[],
): Transition[] {
  return transitionIds.map((transitionId) => {
    const transition = petriNet.transitions.find((t) => t.id === transitionId);
    if (!transition) {
      throw Error(
        `Transition ${transitionId} in petriNet ${petriNet.id} not found`,
      );
    }
    return transition;
  });
}

export function findPlaceMapping(
  petriNet1: PetriNet,
  petriNet2: PetriNet,
): Map<PlaceId, PlaceId> | undefined {
  if (petriNet1.transitions.length !== petriNet2.transitions.length) {
    return undefined;
  }
  const placeMapping = new Map<PlaceId, PlaceId>();
  for (const transition1 of petriNet1.transitions) {
    const transition2 = findTransitions(petriNet2, [transition1.id])[0];
    if (!transition2 || transition1.name !== transition2.name) {
      return undefined;
    }
    if (
      transition1.incomingPlaces.length !== transition2.incomingPlaces.length
    ) {
      return undefined;
    }
    for (const [index, placeId1] of transition1.incomingPlaces.entries()) {
      const placeId2 = placeMapping.get(placeId1);
      if (placeId2 && !transition2.incomingPlaces.includes(placeId2)) {
        return undefined;
      }
      placeMapping.set(placeId1, transition2.incomingPlaces[index]);
    }
    if (
      transition1.outgoingPlaces.length !== transition2.outgoingPlaces.length
    ) {
      return undefined;
    }
    for (const [index, placeId1] of transition1.outgoingPlaces.entries()) {
      const placeId2 = placeMapping.get(placeId1);
      if (placeId2 && !transition2.outgoingPlaces.includes(placeId2)) {
        return undefined;
      }
      placeMapping.set(placeId1, transition2.outgoingPlaces[index]);
    }
  }
  return placeMapping;
}
