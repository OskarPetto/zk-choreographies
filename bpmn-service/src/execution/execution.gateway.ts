import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { Observable, firstValueFrom, map } from 'rxjs';
import { SaltedHash } from '../domain/execution';
import { Model } from 'src/domain/model';

@Injectable()
export class ExecutionGateway {
  constructor(private readonly httpService: HttpService) {}

  async createModel(model: Model): Promise<SaltedHash> {
    const response: Observable<any> = this.httpService.post(
      'http://127.0.0.1:8080/models',
      model,
    );
    return await firstValueFrom(
      response.pipe(map((res: any) => res.data as SaltedHash)),
    );
  }
}
