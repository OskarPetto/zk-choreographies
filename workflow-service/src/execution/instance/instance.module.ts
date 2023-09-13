import { Module } from '@nestjs/common';
import { InstanceService } from './instance.service';

@Module({
  imports: [],
  exports: [],
  controllers: [],
  providers: [InstanceService],
})
export class InstanceModule {}
