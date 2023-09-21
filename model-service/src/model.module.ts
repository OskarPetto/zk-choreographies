import { Module } from '@nestjs/common';
import { BpmnModule } from './bpmn/bpmn.module';
import { PetriNetModule } from './petri-net/petri-net.module';

@Module({
  imports: [BpmnModule, PetriNetModule],
  exports: [],
  controllers: [],
  providers: [],
})
export class ModelModule {}
