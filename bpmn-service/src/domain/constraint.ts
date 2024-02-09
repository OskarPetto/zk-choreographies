import { MessageId } from "./model";

export interface Constraint {
  coefficients: number[];
  messageIds: MessageId[];
  offset: number;
  comparisonOperator: number;
}
