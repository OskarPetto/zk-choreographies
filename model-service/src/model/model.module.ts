import { Module } from '@nestjs/common';
import { ModelReducer } from './model.reducer';
import { ModelService } from './model.service';
import { ModelController } from './model.controller';

@Module({
  imports: [],
  exports: [ModelReducer, ModelService],
  controllers: [ModelController],
  providers: [ModelService, ModelReducer],
})
export class ModelModule {}
