import { Injectable } from '@nestjs/common';
import { ModelId, Model } from './model';

@Injectable()
export class ModelService {
    models: Map<ModelId, Model>;

    find(modelId: ModelId): Model | undefined {
        return this.models.get(modelId);
    }
}
