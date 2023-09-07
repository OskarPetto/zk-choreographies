import { Injectable } from '@nestjs/common';
import { Instance, ExecutionStatus, copyInstance } from 'src/instance';
import { Element, ElementType } from 'src/model';

@Injectable()
export class ExecutionService {
  executeElements(instance: Instance, elements: Element[]): Instance {
    const newInstance = copyInstance(instance);
    for (const element of elements) {
      this.executeElement(newInstance, element);
    }
    return newInstance;
  }

  private executeElement(instance: Instance, element: Element) {
    if (!this.isElementExecutable(instance, element)) {
      throw Error(`Element ${element.id} is not executable`);
    }
    for (const incomingFlowId of element.incomingFlows) {
      instance.executionStatuses.set(
        incomingFlowId,
        ExecutionStatus.NOT_ACTIVE,
      );
    }
    for (const outgoingFlowId of element.outgoingFlows) {
      instance.executionStatuses.set(outgoingFlowId, ExecutionStatus.ACTIVE);
    }
    if (element.type == ElementType.END) {
      instance.finished = true;
    }
  }

  private isElementExecutable(instance: Instance, element: Element) {
    if (instance.finished) {
      return false;
    }
    return [...element.incomingFlows]
      .map((flowId) => instance.executionStatuses.get(flowId))
      .every((executionStatus) => executionStatus === ExecutionStatus.ACTIVE);
  }
}
