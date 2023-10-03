import { SequenceFlow } from 'src/choreography/choreography';
import { ConstraintParser } from './constraint.parser';
import { defaultConstraint } from './constraint';

describe('ConstraintParser', () => {
  let constraintParser: ConstraintParser;
  const sequenceFlows: SequenceFlow[] = [
    {
      id: 'Flow0',
    },
    {
      id: 'Flow1',
      name: 'product not in stock',
    },
    {
      id: 'Flow2',
      name: 'stock < ordered',
    },
    {
      id: 'Flow3',
      name: '5*a > 2*b - 2',
    },
    {
      id: 'Flow4',
      name: '1+3*y == -1*x',
    },
  ];
  const expected = new Map([
    ['Flow0', defaultConstraint()],
    ['Flow1', defaultConstraint()],
    [
      'Flow2',
      {
        a: 1,
        x: 0,
        b: -1,
        y: 1,
        c: 0,
        comparisonOperator: 2,
      },
    ],
    [
      'Flow3',
      {
        a: 5,
        x: 0,
        b: -2,
        y: 1,
        c: 2,
        comparisonOperator: 1,
      },
    ],
    [
      'Flow4',
      {
        a: 3,
        x: 0,
        b: 1,
        y: 1,
        c: 1,
        comparisonOperator: 0,
      },
    ],
  ]);

  beforeAll(() => {
    constraintParser = new ConstraintParser();
  });

  describe('parseConstraints', () => {
    it('should parse sequenceFlow constraints', () => {
      const result = constraintParser.parseConstraints(sequenceFlows);
      expect(result).toEqual(expected);
    });
  });
});
