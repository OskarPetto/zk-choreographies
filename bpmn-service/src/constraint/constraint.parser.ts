import { Injectable } from '@nestjs/common';
import { Constraint, defaultConstraint } from './constraint';
import { SequenceFlow, SequenceFlowId } from 'src/choreography/choreography';
/* eslint @typescript-eslint/no-var-requires: "off" */
const esprima = require('esprima');

interface ConstraintPart {
  literal: number;
  variable?: number;
}

@Injectable()
export class ConstraintParser {
  parseConstraints(
    sequenceFlows: SequenceFlow[],
  ): Map<SequenceFlowId, Constraint> {
    const constraints = new Map();
    const variables = new Map();
    for (const sequenceFlow of sequenceFlows) {
      if (!sequenceFlow.name) {
        constraints.set(sequenceFlow.id, defaultConstraint());
        continue;
      }
      try {
        const script = esprima.parseScript(sequenceFlow.name);
        if (script.body.length != 1) {
          throw Error('not implemented');
        }
        const expression = script.body[0].expression;
        const constraint = this.parseBooleanExpression(expression, variables);
        constraints.set(sequenceFlow.id, constraint);
      } catch (e) {
        console.log('no valid constraint: ', sequenceFlow.name);
        constraints.set(sequenceFlow.id, defaultConstraint());
      }
    }
    return constraints;
  }

  private parseBooleanExpression(
    expression: any,
    variables: Map<string, number>,
  ): Constraint {
    if (expression.type !== 'BinaryExpression') {
      throw Error('not implemented');
    }
    const comparisonOperator = this.parseComparisonOperator(
      expression.operator,
    );
    const leftParts: ConstraintPart[] = this.parseExpression(
      expression.left,
      variables,
      1,
    );
    const rightParts: ConstraintPart[] = this.parseExpression(
      expression.right,
      variables,
      -1,
    );
    const combinedParts = leftParts.concat(rightParts);
    const resultingParts: ConstraintPart[] = [];
    let c = 0;
    for (const part of combinedParts) {
      if (part.variable !== undefined) {
        resultingParts.push(part);
      } else {
        c += part.literal;
      }
    }
    resultingParts.push({ literal: c });
    if (combinedParts.length > 3) {
      throw Error('not implemented');
    }
    return {
      a: resultingParts[0].literal,
      x: resultingParts[0].variable!,
      b: resultingParts.length > 1 ? resultingParts[1].literal : 0,
      y: resultingParts.length > 1 ? resultingParts[1].variable! : 0,
      c: resultingParts.length > 2 ? resultingParts[2].literal : 0,
      comparisonOperator: comparisonOperator,
    };
  }

  private parseExpression(
    expression: any,
    variables: Map<string, number>,
    factor: number,
  ): ConstraintPart[] {
    if (expression.type == 'Literal') {
      return [
        {
          literal: expression.value * factor,
        },
      ];
    } else if (expression.type == 'Identifier') {
      let variable;
      if (variables.has(expression.name)) {
        variable = variables.get(expression.name);
      } else {
        variable = variables.size;
        variables.set(expression.name, variable);
      }
      return [
        {
          literal: factor,
          variable: variable,
        },
      ];
    } else if (expression.type == 'UnaryExpression') {
      if (expression.operator == '-') {
        return this.parseExpression(
          expression.argument,
          variables,
          factor * -1,
        );
      } else {
        throw Error('not implemented');
      }
    } else if (expression.type == 'BinaryExpression') {
      if (expression.operator == '*') {
        const leftParts = this.parseExpression(
          expression.left,
          variables,
          factor,
        );
        const rightParts = this.parseExpression(
          expression.right,
          variables,
          factor,
        );
        if (leftParts.length > 1 || rightParts.length > 1) {
          throw Error('not implemented');
        }
        if (
          leftParts[0].variable !== undefined &&
          rightParts[0].variable !== undefined
        ) {
          throw Error('not implemented');
        }
        if (
          leftParts[0].variable === undefined &&
          rightParts[0].variable === undefined
        ) {
          throw Error('not implemented');
        }
        const literal =
          leftParts[0].variable === undefined
            ? leftParts[0].literal
            : rightParts[0].literal;
        const variable =
          leftParts[0].variable !== undefined
            ? leftParts[0].variable
            : rightParts[0].variable;
        return [
          {
            literal,
            variable,
          },
        ];
      } else if (expression.operator == '+') {
        const leftParts = this.parseExpression(
          expression.left,
          variables,
          factor,
        );
        const rightParts = this.parseExpression(
          expression.right,
          variables,
          factor,
        );
        return leftParts.concat(rightParts);
      } else if (expression.operator == '-') {
        const leftParts = this.parseExpression(
          expression.left,
          variables,
          factor,
        );
        const rightParts = this.parseExpression(
          expression.right,
          variables,
          -1 * factor,
        );
        return leftParts.concat(rightParts);
      } else {
        throw Error('not implemented');
      }
    } else {
      throw Error('not implemented');
    }
  }

  private parseComparisonOperator(operator: string): number {
    switch (operator) {
      case '==':
        return 0;
      case '>':
        return 1;
      case '<':
        return 2;
      case '>=':
        return 3;
      case '<=':
        return 4;
    }
    throw Error('not implemented');
  }
}
