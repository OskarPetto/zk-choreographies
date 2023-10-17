package testdata

import (
	"execution-service/domain"
	"execution-service/parameters"
	"time"
)

type State struct {
	Instance           domain.Instance
	Model              domain.Model
	SenderSignature    domain.Signature
	RecipientSignature *domain.Signature
	Transition         domain.Transition
	Identity           domain.IdentityId
	ConstraintInput    domain.ConstraintInput
	Message            *domain.Message
}

func GetModel2States(signatureParameters parameters.SignatureParameters) []State {
	model2 := GetModel2()
	order := domain.NewIntegerMessage(2)
	stock := domain.NewIntegerMessage(20)
	confirm := domain.NewBytesMessage([]byte("confirm"))
	invoice := domain.NewBytesMessage([]byte("invoice"))
	shippingAddress := domain.NewBytesMessage([]byte("shipping_address"))
	product := domain.NewBytesMessage([]byte("product"))
	payment := domain.NewBytesMessage([]byte("payment"))

	identity0 := domain.IdentityId(0)
	identity1 := domain.IdentityId(1)
	return []State{
		getModelState(
			signatureParameters,
			model2,
			domain.MaxTransitionCount,
			[]domain.PlaceId{10},
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
			0,
			nil,
			nil,
		),
		getModelState(
			signatureParameters,
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
			0,
			nil,
			nil,
		),
		getModelState(
			signatureParameters,
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
				order.Hash.Hash,
			},
			domain.EmptyConstraintInput(),
			0,
			&identity1,
			&order,
		),
		getModelState(
			signatureParameters,
			model2,
			3,
			[]domain.PlaceId{1},
			[]domain.Hash{
				stock.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash.Hash,
			},
			domain.EmptyConstraintInput(),
			1,
			&identity0,
			&stock,
		),
		getModelState(
			signatureParameters,
			model2,
			11,
			[]domain.PlaceId{2, 3},
			[]domain.Hash{
				stock.Hash.Hash,
				confirm.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash.Hash,
			},
			domain.ConstraintInput{
				Messages: []domain.Message{
					order, stock,
				},
			},
			1,
			&identity0,
			&confirm,
		),
		getModelState(
			signatureParameters,
			model2,
			9,
			[]domain.PlaceId{2, 9},
			[]domain.Hash{
				stock.Hash.Hash,
				confirm.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				invoice.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash.Hash,
			},
			domain.EmptyConstraintInput(),
			1,
			&identity0,
			&invoice,
		),
		getModelState(
			signatureParameters,
			model2,
			10,
			[]domain.PlaceId{2, 5},
			[]domain.Hash{
				stock.Hash.Hash,
				confirm.Hash.Hash,
				payment.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				invoice.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash.Hash,
			},
			domain.EmptyConstraintInput(),
			0,
			&identity1,
			&payment,
		),
		getModelState(
			signatureParameters,
			model2,
			7,
			[]domain.PlaceId{8, 5},
			[]domain.Hash{
				stock.Hash.Hash,
				confirm.Hash.Hash,
				payment.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				invoice.Hash.Hash,
				shippingAddress.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash.Hash,
			},
			domain.EmptyConstraintInput(),
			0,
			&identity1,
			&shippingAddress,
		),
		getModelState(
			signatureParameters,
			model2,
			8,
			[]domain.PlaceId{4, 5},
			[]domain.Hash{
				stock.Hash.Hash,
				confirm.Hash.Hash,
				payment.Hash.Hash,
				product.Hash.Hash,
				domain.EmptyHash(),
				invoice.Hash.Hash,
				shippingAddress.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash.Hash,
			},
			domain.EmptyConstraintInput(),
			1,
			&identity0,
			&product,
		),
		getModelState(
			signatureParameters,
			model2,
			1,
			[]domain.PlaceId{11},
			[]domain.Hash{
				stock.Hash.Hash,
				confirm.Hash.Hash,
				payment.Hash.Hash,
				product.Hash.Hash,
				domain.EmptyHash(),
				invoice.Hash.Hash,
				shippingAddress.Hash.Hash,
				domain.EmptyHash(),
				domain.EmptyHash(),
				order.Hash.Hash,
			},
			domain.EmptyConstraintInput(),
			0,
			nil,
			nil,
		),
	}
}

func getModelState(signatureParameters parameters.SignatureParameters, model domain.Model, transitionIndex uint, activePlaces []domain.PlaceId, hashes []domain.Hash, constraintInput domain.ConstraintInput, sender domain.IdentityId, recipient *domain.IdentityId, message *domain.Message) State {
	tokenCounts := make([]int8, model.PlaceCount)
	for _, placeId := range activePlaces {
		tokenCounts[placeId] = 1
	}
	publicKeys := signatureParameters.GetPublicKeys(int(model.ParticipantCount))
	messageHashes := make([]domain.Hash, model.MessageCount)
	copy(messageHashes, hashes)
	instance := domain.Instance{
		Model:         model.Hash.Hash,
		TokenCounts:   tokenCounts,
		PublicKeys:    publicKeys,
		MessageHashes: messageHashes,
		CreatedAt:     time.Now().Unix(),
	}
	instance.UpdateHash()
	senderPrivateKey := signatureParameters.GetPrivateKeyForIdentity(sender)
	senderSignature := instance.Sign(senderPrivateKey)
	var recipientSignature *domain.Signature
	if recipient != nil {
		recipientPrivateKey := signatureParameters.GetPrivateKeyForIdentity(*recipient)
		tmp := instance.Sign(recipientPrivateKey)
		recipientSignature = &tmp
	}

	var transition domain.Transition
	if transitionIndex < domain.MaxTransitionCount {
		transition = model.Transitions[transitionIndex]
	} else {
		transition = domain.OutOfBoundsTransition()
	}
	return State{
		Model:              model,
		Instance:           instance,
		Transition:         transition,
		SenderSignature:    senderSignature,
		RecipientSignature: recipientSignature,
		Identity:           sender,
		ConstraintInput:    constraintInput,
		Message:            message,
	}
}
