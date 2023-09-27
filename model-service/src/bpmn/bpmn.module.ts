import { Module } from '@nestjs/common';
import { BpmnMapper } from './bpmn.mapper';
import { BpmnParser } from './bpmn.parser';
import { BpmnService } from './bpmn.service';
import { ModelModule } from '../model/model.module';
import { BpmnController } from './bpmn.controller';

@Module({
  imports: [ModelModule],
  controllers: [BpmnController],
  providers: [BpmnMapper, BpmnParser, BpmnService],
})
export class BpmnModule {}
