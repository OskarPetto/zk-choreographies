import { MessageId } from "src/model/model";

export interface Constraint {
  coefficients: number[];
  messageIds: MessageId[];
  offset: number;
  comparisonOperator: number;
}

export function defaultConstraint(): Constraint {
  return {
    coefficients: [],
    messageIds: [],
    offset: 0,
    comparisonOperator: 0,
  };
}
