export interface Constraint {
  a: number;
  x: number;
  b: number;
  y: number;
  c: number;
  comparisonOperator: number;
}

export function defaultConstraint(): Constraint {
  return {
    a: 0,
    b: 0,
    x: 0,
    y: 0,
    c: 0,
    comparisonOperator: 0,
  };
}
