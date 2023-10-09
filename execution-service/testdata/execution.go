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
	stock := domain.NewIntegerMessage(20)
	stockHash := stock.Hash
	confirmHash := domain.NewBytesMessage([]byte("confirm")).Hash
	invoiceHash := domain.NewBytesMessage([]byte("invoice")).Hash
	shippingAddressHash := domain.NewBytesMessage([]byte("shipping_address")).Hash
	productHash := domain.NewBytesMessage([]byte("product")).Hash
	paymentHash := domain.NewBytesMessage([]byte("payment")).Hash
	return []State{
		getModelState(
			model2,
			domain.MaxTransitionCount,
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
				stockHash,
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
				stockHash,
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
				Messages: []domain.Message{
					order,
					stock,
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
				stockHash,
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
				stockHash,
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
				stockHash,
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
				stockHash,
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
				stockHash,
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
				stockHash,
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

func getModelState(model domain.Model, transitionIndex uint, activePlaces []domain.PlaceId, hashes []domain.Hash, constraintInput domain.ConstraintInput, signatureParameters parameters.SignatureParameters, idendity domain.IdentityId) State {
	tokenCounts := make([]int8, model.PlaceCount)
	for _, placeId := range activePlaces {
		tokenCounts[placeId] = 1
	}
	publicKeys := signatureParameters.GetPublicKeys(int(model.ParticipantCount))
	messageHashes := make([][domain.HashSize]byte, model.MessageCount)
	for i, messageHash := range hashes {
		messageHashes[i] = messageHash.Value
	}
	instance := domain.Instance{
		Model:         "modelHash",
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
		CreatedAt:     time.Now().Unix(),
	}
	instance.ComputeHash()
	privateKey := signatureParameters.GetPrivateKeyForIdentity(uint(idendity))
	signature := instance.Sign(privateKey)

	var transition domain.Transition
	if transitionIndex < domain.MaxTransitionCount {
		transition = model.Transitions[transitionIndex]
	} else {
		transition = domain.OutOfBoundsTransition()
	}
	return State{
		Model:           model,
		Instance:        instance,
		Transition:      transition,
		Signature:       signature,
		Identity:        idendity,
		ConstraintInput: constraintInput,
	}
}
