import { PetriNetId } from 'src/model/petri-net/petri-net';

export type InstanceId = string;

export interface Instance {
  id?: InstanceId;
  petriNet: PetriNetId;
  tokenCounts: number[];
}
