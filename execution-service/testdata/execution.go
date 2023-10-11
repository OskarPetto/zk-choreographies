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
	Message         *domain.Message
}

func GetModel2States(signatureParameters parameters.SignatureParameters) []State {
	model2 := GetModel2()
	order := domain.NewMessage(nil, 2)
	stock := domain.NewMessage(nil, 20)
	confirm := domain.NewMessage([]byte("confirm"), 0)
	invoice := domain.NewMessage([]byte("invoice"), 0)
	shippingAddress := domain.NewMessage([]byte("shipping_address"), 0)
	product := domain.NewMessage([]byte("product"), 0)
	payment := domain.NewMessage([]byte("payment"), 0)
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
			nil,
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
			nil,
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
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
			&order,
		),
		getModelState(
			model2,
			3,
			[]domain.PlaceId{1},
			[]domain.Hash{
				stock.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			1,
			nil,
		),
		getModelState(
			model2,
			12,
			[]domain.PlaceId{11},
			[]domain.Hash{
				stock.Hash,
				confirm.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.ConstraintInput{
				Messages: []domain.Message{
					order,
					stock,
				},
			},
			signatureParameters,
			1,
			&confirm,
		),
		getModelState(
			model2,
			13,
			[]domain.PlaceId{2, 3},
			[]domain.Hash{
				stock.Hash,
				confirm.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
			nil,
		),
		getModelState(
			model2,
			10,
			[]domain.PlaceId{2, 10},
			[]domain.Hash{
				stock.Hash,
				confirm.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				invoice.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			1,
			&invoice,
		),
		getModelState(
			model2,
			8,
			[]domain.PlaceId{9, 10},
			[]domain.Hash{
				stock.Hash,
				confirm.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				invoice.Hash,
				shippingAddress.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
			&shippingAddress,
		),
		getModelState(
			model2,
			9,
			[]domain.PlaceId{10, 4},
			[]domain.Hash{
				stock.Hash,
				confirm.Hash,
				domain.EmptyHash(),
				product.Hash,
				domain.EmptyHash(),
				invoice.Hash,
				shippingAddress.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			1,
			&product,
		),
		getModelState(
			model2,
			11,
			[]domain.PlaceId{4, 5},
			[]domain.Hash{
				stock.Hash,
				confirm.Hash,
				payment.Hash,
				product.Hash,
				domain.EmptyHash(),
				invoice.Hash,
				shippingAddress.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
			&payment,
		),
		getModelState(
			model2,
			1,
			[]domain.PlaceId{13},
			[]domain.Hash{
				stock.Hash,
				confirm.Hash,
				payment.Hash,
				product.Hash,
				domain.EmptyHash(),
				invoice.Hash,
				shippingAddress.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash,
			},
			domain.EmptyConstraintInput(),
			signatureParameters,
			0,
			nil,
		),
	}
}

func getModelState(model domain.Model, transitionIndex uint, activePlaces []domain.PlaceId, hashes []domain.Hash, constraintInput domain.ConstraintInput, signatureParameters parameters.SignatureParameters, idendity domain.IdentityId, message *domain.Message) State {
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
		Message:         message,
	}
}
