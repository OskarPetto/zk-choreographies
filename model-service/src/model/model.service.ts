import { Injectable } from '@nestjs/common';
import { ModelId, Model } from './model';
import { v4 as uuid } from 'uuid';

@Injectable()
export class ModelService {
  models: Map<ModelId, Model> = new Map();

  saveModel(model: Model) {
    if (!model.id) {
      model.id = this.createModelId();
      model.createdAt = new Date();
    }
    this.models.set(model.id, model);
  }

  findModelById(modelId: ModelId): Model | undefined {
    const model = this.models.get(modelId);
    return model;
  }

  findAllModels(): Model[] {
    return [...this.models.values()].sort(
      (m1, m2) => m1.createdAt.getTime() - m2.createdAt.getTime(),
    );
  }

  private createModelId(): ModelId {
    return uuid();
  }
}
