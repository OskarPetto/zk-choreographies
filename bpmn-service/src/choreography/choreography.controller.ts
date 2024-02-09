import { Controller, Post, Body, Get } from '@nestjs/common';
import { ChoreographyService } from './choreography.service';
import { Choreography } from 'src/domain/choreography';

export class TransformChoreographyCommand {
  xmlString: string;
}

@Controller('choreographies')
export class ChoreographyController {

  constructor(private choreographyService: ChoreographyService) { }

  @Post()
  async transformChoreography(
    @Body() cmd: TransformChoreographyCommand,
  ): Promise<Choreography> {
    return this.choreographyService.transformChoreography(cmd.xmlString);
  }

  @Get()
  async findAllChoreographies(): Promise<Choreography[]> {
    return this.choreographyService.findAllChoreographies();
  }
}
