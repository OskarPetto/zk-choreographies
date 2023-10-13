import { Controller, Post, Body } from '@nestjs/common';
import { ChoreographyService } from './choreography.service';
import { Model } from 'src/model/model';

export class TransformChoreographyCommand {
  xmlString: string;
}

@Controller('choreographies')
export class ChoreographyController {
  constructor(private choreographyService: ChoreographyService) {}
  @Post()
  async transformChoreography(
    @Body() cmd: TransformChoreographyCommand,
  ): Promise<Model> {
    return this.choreographyService.transformChoreography(cmd.xmlString);
  }
}
