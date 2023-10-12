package domain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"execution-service/utils"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards"
	"github.com/consensys/gnark-crypto/ecc/bn254/twistededwards/eddsa"
)

type State struct {
	Model    Model
	Instance Instance
	Message  *Message
}

func NewState(model Model, instance Instance, message *Message) State {
	return State{
		Model:    model,
		Instance: instance,
		Message:  message,
	}
}

type SerializedState struct {
	Value []byte
}

type EncryptedState struct {
	Value     []byte
	Sender    PublicKey
	Recipient PublicKey
}

func (encryptedState *EncryptedState) Decrypt(privateKey *eddsa.PrivateKey) (SerializedState, error) {
	secretKey := ecdh(privateKey, encryptedState.Sender)

	aes, err := aes.NewCipher(secretKey)
	utils.PanicOnError(err)

	gcm, err := cipher.NewGCM(aes)
	utils.PanicOnError(err)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encryptedState.Value[:nonceSize], encryptedState.Value[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return SerializedState{}, err
	}

	if err != nil {
		return SerializedState{}, err
	}
	return SerializedState{
		Value: plaintext,
	}, nil
}

func (state *SerializedState) Encrypt(sender *eddsa.PrivateKey, recipient PublicKey) EncryptedState {
	secretKey := ecdh(sender, recipient)

	aes, err := aes.NewCipher(secretKey)
	utils.PanicOnError(err)

	gcm, err := cipher.NewGCM(aes)
	utils.PanicOnError(err)

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	utils.PanicOnError(err)

	ciphertext := gcm.Seal(nonce, nonce, state.Value, nil)
	return EncryptedState{
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
