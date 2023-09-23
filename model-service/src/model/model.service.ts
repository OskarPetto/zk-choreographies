import { Injectable } from '@nestjs/common';
import { ModelId, Model } from './model';
import { v4 as uuid } from 'uuid';

@Injectable()
export class ModelService {
  models: Map<ModelId, Model> = new Map();

  saveModel(model: Model) {
    if (!model.id) {
      model.id = this.createModelId();
    }
    this.models.set(model.id, model);
  }

  findModel(modelId: ModelId): Model {
    const model = this.models.get(modelId);
    if (!model) {
      throw Error(`Model ${modelId} not found`);
    }
    return model;
  }

  private createModelId(): ModelId {
    return uuid();
  }
}
