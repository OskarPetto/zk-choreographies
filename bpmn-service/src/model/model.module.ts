import { Module } from '@nestjs/common';
import { ModelReducer } from './model.reducer';

@Module({
  imports: [],
  exports: [ModelReducer],
  controllers: [],
  providers: [ModelReducer],
})
export class ModelModule {}
