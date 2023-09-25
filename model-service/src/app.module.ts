import { Module } from '@nestjs/common';
import { BpmnModule } from './bpmn/bpmn.module';
import { ModelModule } from './model/model.module';

@Module({
  imports: [BpmnModule, ModelModule],
  exports: [],
  controllers: [],
  providers: [],
})
export class AppModule {}
