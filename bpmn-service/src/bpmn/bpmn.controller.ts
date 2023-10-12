import { Body, Controller, Post, Get } from '@nestjs/common';
import { BpmnService } from './bpmn.service';
import { BpmnModel } from './bpmn';

interface CreateBpmnModelCommand {
  xmlString: string;
}

@Controller('bpmn')
export class BpmnController {
  constructor(private bpmnService: BpmnService) {}
  @Post()
  async createBpmnModel(@Body() cmd: CreateBpmnModelCommand): Promise<string> {
    return this.bpmnService.createBpmnModel(cmd.xmlString);
  }

  @Get()
  async findAllBpmnModels(): Promise<BpmnModel[]> {
    return this.bpmnService.findAllBpmnModels();
  }
}
