package proof

import (
	"bytes"
	"execution-service/domain"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
)

type Proof struct {
	Value [8]*big.Int
	Input []*big.Int
}

func toProof(groth16Proof groth16.Proof, inputs ...domain.Hash) (Proof, error) {

	var proof Proof

	var buf bytes.Buffer
	_, err := groth16Proof.WriteRawTo(&buf)
	if err != nil {
		return Proof{}, err
	}
	proofBytes := buf.Bytes()

	const fpSize = fr.Bytes
	proof.Value[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	proof.Value[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	proof.Value[2] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	proof.Value[3] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	proof.Value[4] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	proof.Value[5] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	proof.Value[6] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	proof.Value[7] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	proof.Input = make([]*big.Int, len(inputs))

	for i, input := range inputs {
		proof.Input[i] = new(big.Int).SetBytes(input.Value[:])
	}

	return proof, nil
}
