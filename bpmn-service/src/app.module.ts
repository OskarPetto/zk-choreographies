import { Module } from '@nestjs/common';
import { ChoreographyModule } from './choreography/choreography.module';
import { ConfigModule } from '@nestjs/config';
import { ConstraintModule } from './constraint/constraint.module';

@Module({
  imports: [ChoreographyModule, ConstraintModule, ConfigModule.forRoot()],
  exports: [],
  controllers: [],
  providers: [],
})
export class AppModule {}
