# ZK-Choreographies

## Execution-Service
To generate proving keys and the verifier smart contracts `cd` to the execution-service and run 
``` 
go test  ./... -v -count=1
```
To recreate them you have to delete the contents of the `files/public` folder within the execution-service
```
rm files/public/*
```
and rerun the tests.

## Bpmn-Service

You can run the tests of the bpmn-service with 
```
npm run test
```

## Smart Contracts

You can run the smart contract tests in the `solidity` folder with
```
truffle test
```