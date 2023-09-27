package proof

import (
	"bytes"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
)

type Proof struct {
	A           [2]*big.Int
	B           [2][2]*big.Int
	C           [2]*big.Int
	PublicInput []*big.Int
}

func newProof(groth16Proof groth16.Proof, fullWitness witness.Witness) (Proof, error) {

	var proof Proof

	var buf bytes.Buffer
	_, err := groth16Proof.WriteRawTo(&buf)
	if err != nil {
		return Proof{}, err
	}
	proofBytes := buf.Bytes()

	// proof.Ar, proof.Bs, proof.Krs
	const fpSize = 4 * 8
	proof.A[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	proof.A[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	proof.B[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	proof.B[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	proof.B[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	proof.B[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	proof.C[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	proof.C[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	publicWitness, err := fullWitness.Public()
	if err != nil {
		return Proof{}, err
	}
	publicWitnessBytes, err := publicWitness.MarshalBinary()
	if err != nil {
		return Proof{}, err
	}

	headerByteCount := 12
	fieldElementCount := (len(publicWitnessBytes) - headerByteCount) / fr.Bytes
	proof.PublicInput = make([]*big.Int, fieldElementCount)
	for i := 0; i < fieldElementCount; i++ {
		var fieldElement [fr.Bytes]byte
		copy(fieldElement[:], publicWitnessBytes[headerByteCount+i*fr.Bytes:headerByteCount+(i+1)*fr.Bytes-1])
		proof.PublicInput[i] = new(big.Int).SetBytes(fieldElement[:])
	}

	return proof, nil
}
