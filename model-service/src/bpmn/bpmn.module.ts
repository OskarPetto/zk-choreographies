import { Module } from '@nestjs/common';
import { BpmnMapper } from './bpmn.mapper';
import { BpmnParser } from './bpmn.parser';
import { BpmnService } from './bpmn.service';
import { ModelModule } from 'src/model/model.module';

@Module({
  imports: [ModelModule],
  controllers: [],
  providers: [BpmnMapper, BpmnParser, BpmnService],
})
export class BpmnModule { }
