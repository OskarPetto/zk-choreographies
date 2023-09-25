import {
  PlaceId,
  Model,
  TransitionId,
  Transition,
} from 'src/model/model';

export function logObject(object: any) {
  console.dir(object, { depth: null });
}

export function findTransitions(
  model: Model,
  transitionIds: TransitionId[],
): Transition[] {
  return transitionIds.map((transitionId) => {
    const transition = model.transitions.find((t) => t.id === transitionId);
    if (!transition) {
      throw Error(
        `Transition ${transitionId} in model ${model.id} not found`,
      );
    }
    return transition;
  });
}

export function findPlaceMapping(
  model1: Model,
  model2: Model,
): Map<PlaceId, PlaceId> | undefined {
  if (model1.transitions.length !== model2.transitions.length) {
    return undefined;
  }
  const placeMapping = new Map<PlaceId, PlaceId>();
  for (const transition1 of model1.transitions) {
    const transition2 = findTransitions(model2, [transition1.id])[0];
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
