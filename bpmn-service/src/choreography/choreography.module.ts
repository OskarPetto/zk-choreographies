import { Module } from '@nestjs/common';
import { ChoreographyMapper } from './choreography.mapper';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyService } from './choreography.service';
import { ModelModule } from '../model/model.module';
import { ChoreographyController } from './choreography.controller';
import { ExecutionModule } from 'src/execution/execution.module';
import { ConditionModule } from 'src/condition/condition.module';

@Module({
  imports: [ModelModule, ConditionModule, ExecutionModule],
  controllers: [ChoreographyController],
  providers: [ChoreographyMapper, ChoreographyParser, ChoreographyService],
})
export class ChoreographyModule {}
