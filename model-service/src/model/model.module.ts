import { Module } from '@nestjs/common';
import { ModelReducer } from './model.reducer';
import { ModelService } from './model.service';

@Module({
  imports: [],
  exports: [ModelReducer, ModelService],
  controllers: [],
  providers: [ModelService, ModelReducer],
})
export class ModelModule {}
