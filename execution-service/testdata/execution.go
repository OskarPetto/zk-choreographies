package testdata

import (
	"execution-service/domain"
	"execution-service/parameters"
	"time"
)

type State struct {
	Instance        domain.Instance
	Model           domain.Model
	Signature       domain.Signature
	Transition      domain.Transition
	Identity        domain.IdentityId
	ConstraintInput domain.ConstraintInput
}

func GetModel2States(signatureParameters parameters.SignatureParameters) []State {
	model2 := GetModel2()
	order := domain.NewIntegerMessage(2)
	orderHash := order.Hash
	confirmHash := domain.NewBytesMessage([]byte("Confirm")).Hash
	invoiceHash := domain.NewBytesMessage([]byte("Invoice")).Hash
	shippingAddressHash := domain.NewBytesMessage([]byte("Shipping address")).Hash
	productHash := domain.NewBytesMessage([]byte("Product")).Hash
	paymentHash := domain.NewBytesMessage([]byte("Payment")).Hash
	return []State{
		getModelState(
			model2,
			0,
			[]domain.PlaceId{12},
			[]domain.Hash{
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
		),
		getModelState(
			model2,
			0,
			[]domain.PlaceId{0},
			[]domain.Hash{
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
		),
		getModelState(
			model2,
			2,
			[]domain.PlaceId{7},
			[]domain.Hash{
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
		),
		getModelState(
			model2,
			3,
			[]domain.PlaceId{1},
			[]domain.Hash{
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			1,
		),
		getModelState(
			model2,
			12,
			[]domain.PlaceId{11},
			[]domain.Hash{
				confirmHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.ConstraintInput{
				Messages: [domain.MaxConstraintMessageCount]domain.Message{
					order,
					domain.EmptyMessage(),
				},
			},
			signatureParameters,
			1,
		),
		getModelState(
			model2,
			13,
			[]domain.PlaceId{2, 3},
			[]domain.Hash{
				confirmHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
		),
		getModelState(
			model2,
			10,
			[]domain.PlaceId{2, 10},
			[]domain.Hash{
				confirmHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				invoiceHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			1,
		),
		getModelState(
			model2,
			8,
			[]domain.PlaceId{9, 10},
			[]domain.Hash{
				confirmHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				invoiceHash,
				shippingAddressHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
		),
		getModelState(
			model2,
			9,
			[]domain.PlaceId{10, 4},
			[]domain.Hash{
				confirmHash,
				domain.EmptyHash(),
				productHash,
				domain.EmptyHash(),
				invoiceHash,
				shippingAddressHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			1,
		),
		getModelState(
			model2,
			11,
			[]domain.PlaceId{4, 5},
			[]domain.Hash{
				confirmHash,
				paymentHash,
				productHash,
				domain.EmptyHash(),
				invoiceHash,
				shippingAddressHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
		),
		getModelState(
			model2,
			1,
			[]domain.PlaceId{13},
			[]domain.Hash{
				confirmHash,
				paymentHash,
				productHash,
				domain.EmptyHash(),
				invoiceHash,
				shippingAddressHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				orderHash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
		),
	}
}

func getModelState(model domain.Model, transitionIndex uint, activePlaces []domain.PlaceId, messageHashes []domain.Hash, constraintInput domain.ConstraintInput, signatureParameters parameters.SignatureParameters, idendity domain.IdentityId) State {
	var tokenCountsFixedSize [domain.MaxPlaceCount]int8
	for _, placeId := range activePlaces {
		tokenCountsFixedSize[placeId] = 1
	}
	for i := model.PlaceCount; i < domain.MaxPlaceCount; i++ {
		tokenCountsFixedSize[i] = domain.OutOfBoundsTokenCount
	}
	publicKeys := signatureParameters.GetPublicKeys(int(model.ParticipantCount))
	var publicKeysFixedSize [domain.MaxParticipantCount]domain.PublicKey
	copy(publicKeysFixedSize[:], publicKeys)
	for i := model.ParticipantCount; i < domain.MaxParticipantCount; i++ {
		publicKeysFixedSize[i] = domain.OutOfBoundsPublicKey()
	}
	var messageHashesFixedSize [domain.MaxMessageCount][domain.HashSize]byte
	for i, messageHash := range messageHashes {
		messageHashesFixedSize[i] = messageHash.Value
	}
	for i := len(messageHashes); i < domain.MaxMessageCount; i++ {
		messageHashesFixedSize[i] = domain.OutOfBoundsHash().Value
	}
	instance := domain.Instance{
		Model:         "modelHash",
		TokenCounts:   tokenCountsFixedSize,
		PublicKeys:    publicKeysFixedSize,
		MessageHashes: messageHashesFixedSize,
		CreatedAt:     time.Now().Unix(),
	}
	instance.ComputeHash()
	privateKey := signatureParameters.GetPrivateKeyForIdentity(uint(idendity))
	signature := instance.Sign(privateKey)
	return State{
		Model:           model,
		Instance:        instance,
		Transition:      model.Transitions[transitionIndex],
		Signature:       signature,
		Identity:        idendity,
		ConstraintInput: constraintInput,
	}
}
