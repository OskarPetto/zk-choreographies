import { Injectable } from '@nestjs/common';
import { BpmnParser } from './bpmn.parser';
import { BpmnMapper } from './bpmn.mapper';
import { Model } from '../domain/model';
import { ModelReducer } from './model.reducer';

@Injectable()
export class BpmnService {
  constructor(
    private bpmnParser: BpmnParser,
    private bpmnMapper: BpmnMapper,
    private modelReducer: ModelReducer,
  ) {}

  parseModel(bpmnString: string): Model {
    const definitions = this.bpmnParser.parseBpmn(bpmnString);
    const model = this.bpmnMapper.toModel(definitions.process);
    return this.modelReducer.reduceModel(model);
  }
}
