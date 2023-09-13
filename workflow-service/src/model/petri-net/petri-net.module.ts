import { Module } from '@nestjs/common';
import { PetriNetReducer } from './peri-net.reducer';
import { PetriNetService } from './petri-net.service';

@Module({
  imports: [],
  exports: [PetriNetReducer, PetriNetService],
  controllers: [],
  providers: [PetriNetService, PetriNetReducer],
})
export class PetriNetModule {}
