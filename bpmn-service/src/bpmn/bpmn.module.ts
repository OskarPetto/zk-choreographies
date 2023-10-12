import { Module } from '@nestjs/common';
import { BpmnService } from 'src/bpmn/bpmn.service';
import { BpmnController } from 'src/bpmn/bpmn.controller';

@Module({
  imports: [],
  controllers: [BpmnController],
  providers: [BpmnService],
  exports: [BpmnService],
})
export class BpmnModule {}
