
export type TransitionId = string;

export enum TransitionType {
  START,
  END,
  TASK,
  XOR_SPLIT,
  XOR_JOIN,
  AND_SPLIT,
  AND_JOIN
}

export type PlaceId = number;

export interface Transition {
  id: TransitionId;
  type: TransitionType;
  fromPlaces: PlaceId[];
  toPlaces: PlaceId[];
}

export type ModelId = string;

export interface Model {
  id: ModelId;
  placeCount: number;
  transitions: Map<TransitionId, Transition>;
}

export function findTransition(model: Model, transitionId: TransitionId): Transition {
  const transition = model.transitions.get(transitionId);
  if (!transition) {
    throw Error(`Transition ${transitionId} in model ${model.id} not found`);
  }
  return transition;
}

export function copyModel(model: Model): Model {
  const newTransitions = [...model.transitions.values()].map(transition => ({
    id: transition.id,
    type: transition.type,
    fromPlaces: [...transition.fromPlaces],
    toPlaces: [...transition.toPlaces]
  }))

  return {
    id: model.id,
    placeCount: model.placeCount,
    transitions: new Map(newTransitions.map(t => [t.id, t]))
  }
}

export function modelEquals(model1: Model, model2: Model): boolean {
  return this.findPlaceMapping(model1, model2) !== undefined;
}

export function findPlaceMapping(model1: Model, model2: Model): Map<PlaceId, PlaceId> | undefined {
  if (model1.placeCount != model2.placeCount) {
    return undefined;
  }
  if (model1.transitions.size != model2.transitions.size) {
    return undefined;
  }
  const placeMapping = new Map<PlaceId, PlaceId>();
  for (const transition1 of model1.transitions.values()) {
    const transition2 = model2.transitions.get(transition1.id);
    if (!transition2) {
      return undefined;
    }
    const allPlaces1 = [...transition1.fromPlaces, ...transition1.toPlaces];
    const allPlaces2 = [...transition2.fromPlaces, ...transition2.toPlaces];
    if (allPlaces1.length != allPlaces2.length) {
      return undefined;
    }
    for (const [index, placeId1] of allPlaces1.entries()) {
      const placeId2 = allPlaces2[index];
      if (placeMapping.has(placeId1) && placeMapping.get(placeId1) != placeId2) {
        return undefined;
      }
      placeMapping.set(placeId1, placeId2);
    }
  }
  return placeMapping;
}