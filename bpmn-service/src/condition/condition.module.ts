import { Module } from '@nestjs/common';
import { ConditionParser } from './condition.parser';

@Module({
  imports: [],
  controllers: [],
  providers: [ConditionParser],
  exports: [ConditionParser],
})
export class ConditionModule {}
