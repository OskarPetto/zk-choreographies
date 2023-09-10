import { PlaceId, Model } from 'src/model';

export function findPlaceMapping(
  model1: Model,
  model2: Model,
): Map<PlaceId, PlaceId> | undefined {
  if (model1.transitions.size !== model2.transitions.size) {
    return undefined;
  }
  const placeMapping = new Map<PlaceId, PlaceId>();
  for (const transition1 of model1.transitions.values()) {
    const transition2 = model2.transitions.get(transition1.id);
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
