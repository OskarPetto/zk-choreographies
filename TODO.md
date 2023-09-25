# Evaluation
 Proving Time vs Model Size compared to TK23

# Differences from TK23
- Interpreted Execution
- Separate Circuits for Initiation, Transition and Termination
- I have no method for sharing state between participants
- Signature and Public Keys are private
- Support for Loops
- Message Passing 1-n 3 Varianten (nur sender, sender mit multiple signatures, multiple tasks und AND gateway)
- warum nur message hashes? weil prüfen des contents mit programierbaren expressions zu aufwändig wäre
- reduction funktioniert nicht wenn split auf join gateway folgt