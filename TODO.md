# Evaluation
 Proving Time vs Parameter Size compared to TK23 (multiply all by factors of 2)

# Differences from TK23
- Interpreted Execution
- Separate Circuits for Initiation, Transition and Termination
- I have no method for sharing state between participants
- Signature and Public Keys are private

# Describe
- Message Passing from sender to receiver with 3 variants (only sender, sender with multiple signatures, sender and receiver)
- why no variables? because checking their content in programmable gateways is too expensive
- reduction does not work when join gateway immediately follows split gateway and there is a task in one branch (loop)
- message hashes should be hiding, to hide message content from other participants