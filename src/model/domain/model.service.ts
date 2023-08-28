import { Injectable } from '@nestjs/common';
import { ModelId, Model, TransitionId, Transition } from './model';

@Injectable()
export class ModelService {
    models: Map<ModelId, Model>

    saveModel(model: Model) {
        this.models.set(model.id, model);
    }

    findModel(modelId: ModelId): Model {
        const model = this.models.get(modelId);
        if (!model) {
            throw Error(`Model ${modelId} not found`);
        }
        return model;
    }
}
