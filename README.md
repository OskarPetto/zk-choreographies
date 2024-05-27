# ZK-Choreographies

## Bpmn-Service

You can install the dependencies with

```
npm i
```

You can start the bpmn-service with 

```
npm run start
```

You can run the tests of the bpmn-service with 
```
npm run test
```

## Execution-Service

You can run the execution-service in its folder with 
```
go build
```
and then 
```
./execution-service
```

You can run the tests of the execution-service with 
``` 
go test  ./... -v
```
The measurenents in the paper for proving time were obtained from the console output of the prover tests.

To recreate proving keys and the verifier smart contracts you have to delete the contents of the `files/public` folder within the execution-service and restart the application or the tests. 

## Smart Contracts

You can run the smart contract tests in the `solidity` folder with
```
truffle test
```
The measurenents in the paper for gas usage were obtained by running 
```
truffle develop
```
and then 
```
migrate
```

## E2E Test 
After starting the bpmn-service and the execution service, you can run the test in the `e2e` folder with 
```
python e2e.py
```
