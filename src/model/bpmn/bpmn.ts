import * as fs from 'fs';
import * as path from 'path';
import BpmnModdle from 'bpmn-moddle';

const moddle = new BpmnModdle();

export async function readBPMN() {

    const filePath = path.join(process.cwd(), 'test/data/conformance_example.bpmn');
    const xmlStr = fs.readFileSync(filePath, 'utf-8').toString();

    const definitions = await moddle.fromXML(xmlStr);

    //console.log(definitions);
}