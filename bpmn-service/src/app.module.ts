import { Module } from '@nestjs/common';
import { ChoreographyModule } from './choreography/choreography.module';
import { ModelModule } from './model/model.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [ChoreographyModule, ModelModule, ConfigModule.forRoot()],
  exports: [],
  controllers: [],
  providers: [],
})
export class AppModule {}
