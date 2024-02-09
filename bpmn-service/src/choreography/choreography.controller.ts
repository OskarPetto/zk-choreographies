import { Controller, Post, Body, Logger, Get } from '@nestjs/common';
import { ChoreographyService } from './choreography.service';
import { Choreography } from 'src/domain/choreography';

export class TransformChoreographyCommand {
  xmlString: string;
}

@Controller('choreographies')
export class ChoreographyController {
  private readonly logger = new Logger(ChoreographyController.name);

  constructor(private choreographyService: ChoreographyService) {}

  @Post()
  async transformChoreography(
    @Body() cmd: TransformChoreographyCommand,
  ): Promise<Choreography> {
    this.logger.log('Received TransformChoreographyCommand');
    return this.choreographyService.transformChoreography(cmd.xmlString);
  }

  @Get()
  async findAllChoreographies(): Promise<Choreography[]> {
    this.logger.log('Find all choreographies');
    return this.choreographyService.findAllChoreographies();
  }
}
