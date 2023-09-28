import { Controller, Post, Body } from '@nestjs/common';
import { ChoreographyService } from './choreography.service';

export class ImportChoreographyCommand {
  bpmnString: string;
}

@Controller('choreographies')
export class ChoreographyController {
  constructor(private choreographyService: ChoreographyService) {}
  @Post()
  importChoreography(@Body() cmd: ImportChoreographyCommand) {
    this.choreographyService.importChoreography(cmd.bpmnString);
  }
}
