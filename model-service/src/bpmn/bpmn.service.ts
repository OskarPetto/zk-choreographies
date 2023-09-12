import { Injectable } from '@nestjs/common';
import { BpmnParser } from './bpmn.parser';
import { BpmnMapper } from './bpmn.mapper';
import { ModelReducer } from 'src/model/model.reducer';
import { ModelService } from 'src/model/model.service';

@Injectable()
export class BpmnService {
  constructor(
    private bpmnParser: BpmnParser,
    private bpmnMapper: BpmnMapper,
    private modelReducer: ModelReducer,
    private modelService: ModelService,
  ) { }

  importBpmn(bpmnString: string) {
    const definitions = this.bpmnParser.parseBpmn(bpmnString);
    const model = this.bpmnMapper.toModel(definitions.process);
    const reducedModel = this.modelReducer.reduceModel(model);
    this.modelService.saveModel(reducedModel);
  }
}
