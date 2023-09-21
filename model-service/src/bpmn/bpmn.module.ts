import { Module } from '@nestjs/common';
import { BpmnMapper } from './bpmn.mapper';
import { BpmnParser } from './bpmn.parser';
import { BpmnService } from './bpmn.service';
import { PetriNetModule } from '../petri-net/petri-net.module';

@Module({
  imports: [PetriNetModule],
  controllers: [],
  providers: [BpmnMapper, BpmnParser, BpmnService],
})
export class BpmnModule {}
