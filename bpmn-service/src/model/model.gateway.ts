import { Injectable } from '@nestjs/common';
import { Model } from './model';
import { HttpService } from '@nestjs/axios';
import { AxiosError, AxiosResponse } from 'axios';
import { Observable, catchError, firstValueFrom, map } from 'rxjs';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class ModelGateway {
  constructor(
    private readonly httpService: HttpService,
    private configService: ConfigService,
  ) { }
  async createModel(model: Model): Promise<Model> {
    const createModelUrl =
      this.configService.get<string>('EXECUTION_SERVICE_URL') + '/models';
    const { data } = await firstValueFrom(
      this.httpService.post(createModelUrl, model).pipe(
        catchError((error: AxiosError) => {
          throw Error(`Model could not be imported because : ${error?.response?.data}`);
        }),
      )
    );
    return data;
  }
}
