import { Injectable } from '@nestjs/common';
import { ModelId, Model } from './model';
import { TransitionId } from 'model';

@Injectable()
export class ModelService {
    models: Map<ModelId, Model>;

    find(modelId: ModelId): Model {
        return this.models.get(modelId);
    }
}
