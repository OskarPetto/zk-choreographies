import { MessageId } from './model';

export interface Condition {
  coefficients: number[];
  messageIds: MessageId[];
  offset: number;
  comparisonOperator: number;
}
