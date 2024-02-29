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
instanceCreatedEvent = response.json()

instance0 = instanceCreatedEvent['instance']
proof0 = instanceCreatedEvent['proof']

executeTransitionCommand = {
    'instance': instance0['id'],
    'transition': model['transitions'][0]['id'],
    'identity': 0
}

response = requests.post('http://localhost:8080/execution/executeTransition', json=executeTransitionCommand)
instanceCreatedEvent = response.json()

instance1 = instanceCreatedEvent['instance']
proof1 = instanceCreatedEvent['proof']

createInitiatingMessageCommand = {
    'instance': instance1['id'],
    'transition': model['transitions'][15]['id'],
    'bytesMessage': base64.b32encode(bytearray("mountain_bike", 'ascii')).decode('utf-8'),
    'identity': 1
}

response = requests.post('http://localhost:8080/execution/createInitiatingMessage', json=createInitiatingMessageCommand)
initiatingMessageCreatedEvent = response.json()

initiatingMessage = initiatingMessageCreatedEvent['initiatingMessage']

receiveInitiatingMessageCommand = {
    'model': model,
    'instance': instance1,
    'transition': model['transitions'][15]['id'],
    'initiatingMessage': initiatingMessage,
    'identity': 0
}

response = requests.post('http://localhost:8080/execution/receiveInitiatingMessage', json=receiveInitiatingMessageCommand)
initiatingMessageReceivedEvent = response.json()

instance2 = initiatingMessageReceivedEvent['nextInstance']
respondingParticipantSignature = initiatingMessageReceivedEvent['respondingParticipantSignature']

proveMessageExchangeCommand = {
    'currentInstance': instance1['id'],
    'transition': model['transitions'][15]['id'],
    'initiatingMessage': initiatingMessage['id'],
    'identity': 1,
    'nextInstance': instance2,
    'respondingParticipantSignature': respondingParticipantSignature
}

response = requests.post('http://localhost:8080/execution/proveMessageExchange', json=proveMessageExchangeCommand)
instanceCreatedEvent = response.json()

proof2 = instanceCreatedEvent['proof']

fakeTransitionCommand = {
    'instance': instance2['id'],
    'identity': 0,
}

response = requests.post('http://localhost:8080/execution/fakeTransition', json=fakeTransitionCommand)
instanceCreatedEvent = response.json()

instance3 = instanceCreatedEvent['instance']
proof3 = instanceCreatedEvent['proof']

print(instance3)