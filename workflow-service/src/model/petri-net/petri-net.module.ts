import { Module } from '@nestjs/common';
import { PetriNetReducer } from './perti-net.reducer';
import { PetriNetService } from './petri-net.service';

@Module({
  imports: [],
  exports: [PetriNetReducer, PetriNetService],
  controllers: [],
  providers: [PetriNetService, PetriNetReducer],
})
export class PetriNetModule { }
