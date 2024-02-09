import { HttpService } from "@nestjs/axios";
import { Injectable } from "@nestjs/common";
import { AxiosResponse } from "axios";
import { Observable, firstValueFrom } from "rxjs";
import { SaltedHash } from "../domain/execution";
import { Model } from "src/domain/model";

@Injectable()
export class ExecutionGateway {
  constructor(private readonly httpService: HttpService) { }

  async createModel(model: Model): Promise<SaltedHash> {
    return await firstValueFrom(this.httpService.post('http://localhost:8080/models', model)).data;
  }
}
