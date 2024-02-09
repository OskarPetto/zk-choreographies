import { Module } from '@nestjs/common';
import { ExecutionGateway } from './execution.gateway';
import { HttpModule } from '@nestjs/axios';

@Module({
  imports: [HttpModule],
  controllers: [],
  providers: [ExecutionGateway],
  exports: [ExecutionGateway],
})
export class ExecutionModule {}
