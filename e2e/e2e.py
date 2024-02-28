import requests
import base64

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


response = requests.post('http://localhost:8080/execution/executeTransition', json=executeTransitionCommand)
executedTransitionEvent = response.json()

instance2 = executedTransitionEvent['instance']
proof2 = executedTransitionEvent['proof']

createInitiatingMessageCommand = {
    'model': modelId,
    'instance': instance2['id'],
    'transition': model['transitions'][15]['id'],
    'bytesMessage': base64.b32encode(bytearray("mountain_bike", 'ascii')).decode('utf-8'),
    'identity': 1
}

response = requests.post('http://localhost:8080/execution/createInitiatingMessage', json=createInitiatingMessageCommand)
createdInitiatingMessageEvent = response.json()

initiatingMessage = createdInitiatingMessageEvent['initiatingMessage']

receiveInitiatingMessageCommand = {
    'model': model,
    'instance': instance2,
    'transition': model['transitions'][15]['id'],
    'initiatingMessage': initiatingMessage,
    'identity': 0
}

response = requests.post('http://localhost:8080/execution/receiveInitiatingMessage', json=receiveInitiatingMessageCommand)
receivedInitiatingMessageEvent = response.json()

instance3 = receivedInitiatingMessageEvent['nextInstance']
respondingParticipantSignature = receivedInitiatingMessageEvent['respondingParticipantSignature']

proveMessageExchangeCommand = {
    'model': modelId,
    'currentInstance': instance2['id'],
    'transition': model['transitions'][15]['id'],
    'initiatingMessage': initiatingMessage['id'],
    'identity': 1,
    'nextInstance': instance3,
    'respondingParticipantSignature': respondingParticipantSignature
}

response = requests.post('http://localhost:8080/execution/proveMessageExchange', json=proveMessageExchangeCommand)
provedMessageExchangeEvent = response.json()

print(provedMessageExchangeEvent)