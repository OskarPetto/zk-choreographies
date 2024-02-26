import requests
import json

filename = '../bpmn-service/test/data/BikeRental_example.bpmn'
bpmnString = ''

with open(filename, 'r') as file:
    bpmnString = file.read()

transformChoreographyCommand = {'xmlString': bpmnString}
response = requests.post('http://localhost:3000/choreographies', json=transformChoreographyCommand)
modelId = response.json()['id']

response = requests.get('http://localhost:8080/models/' + modelId)
model = response.json()

response = requests.get('http://localhost:8080/publicKeys')
publicKeys = response.json()

instantiateModelCommand = {
    'model': modelId,
    'publicKeys': publicKeys,
    'identity': 0
}

response = requests.post('http://localhost:8080/execution/instantiateModel', json=instantiateModelCommand)
instantiatedModelEvent = response.json()

instance0 = instantiatedModelEvent['instance']
proof0 = instantiatedModelEvent['proof']

executeTransitionCommand = {
    'model': modelId,
    'instance': instance0['id'],
    'transition': model['transitions'][0]['id'],
    'identity': 0
}

response = requests.post('http://localhost:8080/execution/executeTransition', json=executeTransitionCommand)
executedTransitionEvent = response.json()

instance1 = executedTransitionEvent['instance']
proof1 = executedTransitionEvent['proof']

print(instance1)