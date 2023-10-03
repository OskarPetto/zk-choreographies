import { Module } from '@nestjs/common';
import { ConstraintParser } from './constraint.parser';

@Module({
  imports: [],
  controllers: [],
  providers: [ConstraintParser],
})
export class ConstraintModule {}
