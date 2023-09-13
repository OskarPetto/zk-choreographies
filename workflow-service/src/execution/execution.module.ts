import { Module } from '@nestjs/common';
import { ExecutionService } from './execution.service';
import { InstanceModule } from 'src/execution/instance/instance.module';

@Module({
  imports: [InstanceModule],
  exports: [],
  controllers: [],
  providers: [ExecutionService],
})
export class ExecutionModule {}
