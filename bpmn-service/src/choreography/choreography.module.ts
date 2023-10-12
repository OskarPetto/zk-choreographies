import { Module } from '@nestjs/common';
import { ChoreographyMapper } from './choreography.mapper';
import { ChoreographyParser } from './choreography.parser';
import { ChoreographyService } from './choreography.service';
import { ModelModule } from '../model/model.module';
import { ChoreographyController } from './choreography.controller';
import { BpmnModule } from 'src/bpmn/bpmn.module';
import { ConstraintModule } from 'src/constraint/constraint.module';

@Module({
  imports: [ModelModule, BpmnModule, ConstraintModule],
  controllers: [ChoreographyController],
  providers: [ChoreographyMapper, ChoreographyParser, ChoreographyService],
})
export class ChoreographyModule {}
