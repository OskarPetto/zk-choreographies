import { HttpService } from '@nestjs/axios';
import { Injectable } from '@nestjs/common';
import { Observable, catchError, firstValueFrom, map, throwError } from 'rxjs';
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
      response.pipe(
        map((res: any) => res.data as SaltedHash),
        catchError((error: any) => {
          console.error('Error creating model:', error.response.config.data);
          return throwError('Failed to create model. Please try again.'); // You can customize this error message
        }),
      ),
    );
  }
}
