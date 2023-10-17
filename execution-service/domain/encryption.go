package domain

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"execution-service/utils"
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

type Plaintext struct {
	Value []byte
}

type Ciphertext struct {
	Value     []byte
	Sender    PublicKey
	Recipient PublicKey
}

func (ciphertext *Ciphertext) Decrypt(privateKey *eddsa.PrivateKey) (Plaintext, error) {
	publicKey := privateKey.PublicKey.A.Bytes()
	if !bytes.Equal(publicKey[:], ciphertext.Recipient.Value) {
		return Plaintext{}, fmt.Errorf("ciphertext is meant for %s", utils.BytesToString(publicKey[:]))
	}
	secretKey := ecdh(privateKey, ciphertext.Sender)

	aes, err := aes.NewCipher(secretKey)
	utils.PanicOnError(err)

	gcm, err := cipher.NewGCM(aes)
	utils.PanicOnError(err)

	nonceSize := gcm.NonceSize()
	nonce, remainingCiphertext := ciphertext.Value[:nonceSize], ciphertext.Value[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, remainingCiphertext, nil)
	if err != nil {
		return Plaintext{}, err
	}

	if err != nil {
		return Plaintext{}, err
	}
	return Plaintext{
		Value: plaintext,
	}, nil
}

func (plaintext *Plaintext) Encrypt(sender *eddsa.PrivateKey, recipient PublicKey) Ciphertext {
	secretKey := ecdh(sender, recipient)

	aes, err := aes.NewCipher(secretKey)
	utils.PanicOnError(err)

	gcm, err := cipher.NewGCM(aes)
	utils.PanicOnError(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	utils.PanicOnError(err)

	ciphertext := gcm.Seal(nonce, nonce, plaintext.Value, nil)
	return Ciphertext{
		Value:     ciphertext,
		Sender:    NewPublicKey(sender.PublicKey),
		Recipient: recipient,
	}
}

func ecdh(privateKey *eddsa.PrivateKey, publicKey PublicKey) []byte {
	privateKeyBytes := privateKey.Bytes()
	scalarBytes := privateKeyBytes[fr.Bytes : 2*fr.Bytes]
	scalar := new(big.Int).SetBytes(scalarBytes)
	var eddsaPublicKey eddsa.PublicKey
	eddsaPublicKey.A.SetBytes(publicKey.Value)
	var sharedSecret twistededwards.PointAffine
	sharedSecret.ScalarMultiplication(&eddsaPublicKey.A, scalar)
	sharedSecretBytes := sharedSecret.X.Bytes()
	secretKey := sha256.Sum256(sharedSecretBytes[:])
	return secretKey[:]
}
