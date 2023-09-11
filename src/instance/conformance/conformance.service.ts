import { Injectable } from '@nestjs/common';
import { Instance } from '../domain/instance';
import { Model } from 'src/model';

@Injectable()
export class ConformanceService {
  isDerivable(
    instanceBefore: Instance,
    instanceAfter: Instance,
    model: Model,
  ): boolean {
    if (!this.isInstanceValid(instanceBefore, model)) {
      return false;
    }
    if (!this.isInstanceValid(instanceAfter, model)) {
      return false;
    }
    return false;
  }

  isInstanceValid(instance: Instance, model: Model): boolean {
    if (!this.isInstanceOfModel(instance, model)) {
      return false;
    }
    return false;
  }

  private isInstanceOfModel(instance: Instance, model: Model): boolean {
    return instance.tokenCounts.length === model.placeCount;
  }
}
