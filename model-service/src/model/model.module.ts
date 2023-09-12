import { Module } from '@nestjs/common';
import { ModelService } from './model.service';
import { ModelReducer } from './model.reducer';

@Module({
  imports: [],
  exports: [ModelReducer, ModelService],
  controllers: [],
  providers: [ModelService, ModelReducer],
})
export class ModelModule { }
