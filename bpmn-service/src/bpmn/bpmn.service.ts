import { Injectable } from '@nestjs/common';
import { BpmnModel } from './bpmn';
import { v4 as uuid } from 'uuid';

@Injectable()
export class BpmnService {
  private bpmnModels: Map<string, BpmnModel> = new Map();

  createBpmnModel(xmlString: string): string {
    const id = uuid();
    const model = {
      id,
      xmlString,
    };
    this.bpmnModels.set(model.id, model);
    return model.id;
  }

  findBpmnModelById(id: string): BpmnModel {
    if (!this.bpmnModels.has(id)) {
      throw Error(`BpmnModel ${id} does not exist`);
    }
    return this.bpmnModels.get(id)!;
  }

  findAllBpmnModels(): BpmnModel[] {
    return [...this.bpmnModels.values()];
  }
}
