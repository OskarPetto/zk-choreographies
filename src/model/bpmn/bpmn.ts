import * as fs from 'fs';
import * as path from 'path';
import { logObject } from 'src/utils/utils';
const BpmnModdle = require('bpmn-moddle');

const moddle = new BpmnModdle();

export async function readBPMN() {

    const filePath = path.join(process.cwd(), 'test/data/conformance_example.bpmn');
    const xmlStr = fs.readFileSync(filePath, 'utf-8').toString();

    const definitions = await moddle.fromXML(xmlStr);

    const outputFilePath = path.join(process.cwd(), 'test/data/conformance_example.json');
    fs.writeFileSync(outputFilePath, JSON.stringify(definitions));
    // const rootElement = definitions.rootElement;
    // console.log(rootElement);
    // console.log(rootElement.rootElements);
    // console.log(rootElement.rootElements[0].flowElements);
    // console.log(rootElement.rootElements[0].flowElements[0].$type);
}