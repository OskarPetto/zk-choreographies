import { Controller, Post, Body } from '@nestjs/common';
import { BpmnService } from 'src/bpmn/bpmn.service';

export class ImportBpmnCommand {
  bpmnString: string;
}

@Controller('bpmn')
export class BpmnController {
  constructor(private bpmnService: BpmnService) {}
  @Post()
  importBpmn(@Body() cmd: ImportBpmnCommand) {
    this.bpmnService.importBpmn(cmd.bpmnString);
  }
}
