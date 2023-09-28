import { Module } from '@nestjs/common';
import { ModelReducer } from './model.reducer';
import { ModelGateway } from './model.gateway';
import { HttpModule } from '@nestjs/axios';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [HttpModule, ConfigModule],
  exports: [ModelReducer, ModelGateway],
  controllers: [],
  providers: [ModelGateway, ModelReducer],
})
export class ModelModule {}
