import { Injectable } from '@nestjs/common';
import { Model } from './model';
import { HttpService } from '@nestjs/axios';
import { AxiosResponse } from 'axios';
import { Observable } from 'rxjs';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ModelGateway {
  constructor(
    private readonly httpService: HttpService,
    private configService: ConfigService,
  ) {}
  createModel(model: Model): Observable<AxiosResponse<unknown>> {
    const createModelUrl =
      this.configService.get<string>('EXECUTION_SERVICE_URL') + '/models';
    return this.httpService.post(createModelUrl, model);
  }
}
