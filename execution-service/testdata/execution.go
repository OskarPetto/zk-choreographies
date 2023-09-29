package testdata

import (
	"execution-service/domain"
	"execution-service/parameters"
	"time"
)

type State struct {
	Instance  domain.Instance
	Model     domain.Model
	Signature domain.Signature
	Identity  domain.IdentityId
}

func GetModel2States(signatureParameters parameters.SignatureParameters) []State {
	model2 := GetModel2()
	purchaseOrderHash := domain.HashMessage([]byte("Purchase order"))
	confirmHash := domain.HashMessage([]byte("Confirm"))
	invoiceHash := domain.HashMessage([]byte("Invoice"))
	shippingAddressHash := domain.HashMessage([]byte("Shipping address"))
	productHash := domain.HashMessage([]byte("Product"))
	paymentHash := domain.HashMessage([]byte("Payment"))
	return []State{
		getModelState(
			model2,
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
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			0,
		),
		getModelState(
			model2,
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
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			0,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			0,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			1,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			1,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			0,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			1,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			0,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			1,
		),
		getModelState(
			model2,
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
				purchaseOrderHash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
			},
			signatureParameters,
			0,
		),
		getModelState(
			model2,
			[]domain.PlaceId{13},
			[]domain.Hash{
				confirmHash,
				paymentHash,
				productHash,
				domain.Hash{},
				invoiceHash,
				shippingAddressHash,
				domain.Hash{},
				domain.Hash{},
				purchaseOrderHash,
				domain.Hash{},
				domain.Hash{},
				domain.Hash{},
			},
			signatureParameters,
			0,
		),
	}
}

func getModelState(model domain.Model, activePlaces []domain.PlaceId, messageHashes []domain.Hash, signatureParameters parameters.SignatureParameters, idendity domain.IdentityId) State {
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
	var messageHashesFixedSize [domain.MaxMessageCount]domain.Hash
	copy(messageHashesFixedSize[:], messageHashes)
	for i := len(messageHashes); i < domain.MaxMessageCount; i++ {
		messageHashesFixedSize[i] = domain.OutOfBoundsHash()
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
		Model:     model,
		Instance:  instance,
		Signature: signature,
		Identity:  idendity,
	}
}
