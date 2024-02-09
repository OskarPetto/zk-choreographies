import { Module } from '@nestjs/common';
import { ChoreographyModule } from './choreography/choreography.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [ChoreographyModule, ConfigModule.forRoot()],
  exports: [],
  controllers: [],
  providers: [],
})
export class AppModule {}
