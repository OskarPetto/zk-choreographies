import { Injectable } from '@nestjs/common';
import { Model, PlaceId, Transition, TransitionType, copyModel } from './model';

@Injectable()
export class ReductionService {
    reduceModel(model: Model): Model {
        const newModel = copyModel(model);
        for (const transition of newModel.transitions.values()) {
            switch (transition.type) {
                case TransitionType.START:
                case TransitionType.END:
                case TransitionType.TASK:
                    break;
                case TransitionType.XOR_SPLIT:
                case TransitionType.AND_JOIN:
                    this.removeTransitionAndToPlaces(newModel, transition);
                    break;
                case TransitionType.XOR_JOIN:
                case TransitionType.AND_SPLIT:
                    this.removeTransitionAndFromPlaces(newModel, transition);
                    break;
            }
        }
        this.updatePlaceCount(newModel)
        return newModel;
    }

    private updatePlaceCount(model: Model) {
        let places: PlaceId[] = [];
        for (const transition of model.transitions.values()) {
            places = this.union(places, [...transition.fromPlaces, ...transition.toPlaces]);
        }
        model.placeCount = places.length;
    }

    private removeTransitionAndToPlaces(model: Model, transitionToRemove: Transition) {
        for (const transition of model.transitions.values()) {
            const intersect = this.intersect(transition.fromPlaces, transitionToRemove.toPlaces);
            if (intersect.length > 0) {
                transition.fromPlaces = this.setMinus(transition.fromPlaces, transitionToRemove.toPlaces);
                transition.fromPlaces = this.union(transition.fromPlaces, transitionToRemove.fromPlaces);
            }
        }
        model.transitions.delete(transitionToRemove.id);
    }

    private removeTransitionAndFromPlaces(model: Model, transitionToRemove: Transition) {
        for (const transition of model.transitions.values()) {
            const intersect = this.intersect(transition.toPlaces, transitionToRemove.fromPlaces);
            if (intersect.length > 0) {
                transition.toPlaces = this.setMinus(transition.toPlaces, transitionToRemove.fromPlaces);
                transition.toPlaces = this.union(transition.toPlaces, transitionToRemove.toPlaces);
            }
        }
        model.transitions.delete(transitionToRemove.id);
    }

    private setMinus(places1: PlaceId[], places2: PlaceId[]): PlaceId[] {
        return places1.filter(place1 => !places2.includes(place1));
    }

    private intersect(places1: PlaceId[], places2: PlaceId[]): PlaceId[] {
        return places1.filter(place1 => places2.includes(place1));
    }

    private union(places1: PlaceId[], places2: PlaceId[]): PlaceId[] {
        return [...new Set([...places1, ...places2])];
    }
}
