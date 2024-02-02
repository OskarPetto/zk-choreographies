import { Controller, Post, Body, Logger } from '@nestjs/common';
import { ChoreographyService } from './choreography.service';
import { Model } from 'src/model/model';

export class TransformChoreographyCommand {
  xmlString: string;
}

@Controller('choreographies')
export class ChoreographyController {
  private readonly logger = new Logger(ChoreographyController.name);

  constructor(private choreographyService: ChoreographyService) { }
  @Post()
  async transformChoreography(
    @Body() cmd: TransformChoreographyCommand,
  ): Promise<Model> {
    this.logger.log('Received TransformChoreographyCommand');
    return this.choreographyService.transformChoreography(cmd.xmlString);
  }
}
