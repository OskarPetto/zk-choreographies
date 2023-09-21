import { Module } from '@nestjs/common';
import { ModelModule } from './model/model.module';
import { ExecutionModule } from './execution/execution.module';

@Module({
  imports: [ModelModule, ExecutionModule],
  controllers: [],
  providers: [],
})
export class AppModule {}
