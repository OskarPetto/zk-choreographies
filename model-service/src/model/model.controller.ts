import { Controller, Get, NotFoundException, Param } from '@nestjs/common';
import { ModelService } from './model.service';
import { Model } from './model';

@Controller('models')
export class ModelController {
  constructor(private modelService: ModelService) {}
  @Get('/:id')
  findModelById(@Param('id') id: string): Model {
    const result = this.modelService.findModelById(id);
    if (!result) {
      throw new NotFoundException(`model ${id} not found`);
    }
    return result;
  }
  @Get()
  findAllModels(): Model[] {
    return this.modelService.findAllModels();
  }
}
