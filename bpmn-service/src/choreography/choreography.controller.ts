import { Controller, Post, Body } from '@nestjs/common';
import { ChoreographyService } from './choreography.service';

export class ImportChoreographyCommand {
  xmlString: string;
}

@Controller('choreographies')
export class ChoreographyController {
  constructor(private choreographyService: ChoreographyService) {}
  @Post()
  async importChoreography(
    @Body() cmd: ImportChoreographyCommand,
  ): Promise<string> {
    return this.choreographyService.transformChoreography(cmd.xmlString);
  }
}
