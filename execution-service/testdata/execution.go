package testdata

import (
	"execution-service/domain"
	"execution-service/parameters"
	"execution-service/utils"
)

type State struct {
	Instance                       domain.Instance
	Model                          domain.Model
	InitiatingParticipantSignature domain.Signature
	RespondingParticipantSignature *domain.Signature
	Transition                     domain.Transition
	InitiatingParticipant          domain.IdentityId
	RespondingParticipant          *domain.IdentityId
	ConstraintInput                domain.ConstraintInput
	InitiatingMessage              *domain.Message
	RespondingMessage              *domain.Message
}

func GetModel2States(signatureParameters parameters.SignatureParameters) []State {
	model2 := GetModel2()

	state0 := instantiateModel(signatureParameters, model2)
	state1 := executeTransition(signatureParameters, model2, 0, state0.Instance, nil, nil, nil)
	order := domain.NewIntegerMessage(state1.Instance, 2)
	stock := domain.NewIntegerMessage(state1.Instance, 20)
	state2 := executeTransition(signatureParameters, model2, 2, state1.Instance, nil, &order, &stock)
	confirm := domain.NewBytesMessage(state2.Instance, []byte("confirm"))
	constraintInput1 := domain.ConstraintInput{
		Messages: []domain.Message{
			order, stock,
		},
	}
	state3 := executeTransition(signatureParameters, model2, 7, state2.Instance, &constraintInput1, &confirm, nil)
	invoice := domain.NewBytesMessage(state3.Instance, []byte("invoice"))
	payment := domain.NewBytesMessage(state3.Instance, []byte("payment"))
	state4 := executeTransition(signatureParameters, model2, 6, state3.Instance, nil, &invoice, &payment)
	shippingAddress := domain.NewBytesMessage(state4.Instance, []byte("shipping_address"))
	product := domain.NewBytesMessage(state4.Instance, []byte("product"))
	state5 := executeTransition(signatureParameters, model2, 5, state4.Instance, nil, &shippingAddress, &product)
	state6 := executeTransition(signatureParameters, model2, 1, state5.Instance, nil, nil, nil)

	return []State{
		state0, state1, state2, state3, state4, state5, state6,
	}
}

func instantiateModel(signatureParameters parameters.SignatureParameters, model domain.Model) State {
	publicKeys := signatureParameters.GetPublicKeys(int(model.ParticipantCount))
	instance, err := model.Instantiate(publicKeys)
	utils.PanicOnError(err)

	initiatingParticipantPrivateKey := signatureParameters.GetPrivateKeyForIdentity(0)
	initiatingParticipantSignature := instance.Sign(initiatingParticipantPrivateKey)
	return State{
		Model:                          model,
		Instance:                       instance,
		Transition:                     domain.OutOfBoundsTransition(),
		InitiatingParticipantSignature: initiatingParticipantSignature,
	}
}

func executeTransition(signatureParameters parameters.SignatureParameters, model domain.Model, transitionIndex uint, currentInstance domain.Instance, constraintInput *domain.ConstraintInput, initiatingMessage *domain.Message, respondingMessage *domain.Message) State {

	var transition domain.Transition
	if transitionIndex < domain.MaxTransitionCount {
		transition = model.Transitions[transitionIndex]
	} else {
		transition = domain.OutOfBoundsTransition()
	}
	constraintInputNotNull := domain.EmptyConstraintInput()
	if constraintInput != nil {
		constraintInputNotNull = *constraintInput
	}

	nextInstance, err := currentInstance.ExecuteTransition(transition, constraintInputNotNull, initiatingMessage, respondingMessage)
	utils.PanicOnError(err)

	initiatingParticipant := uint(0)
	if transition.InitiatingParticipant != domain.EmptyParticipantId {
		initiatingParticipant = uint(transition.InitiatingParticipant)
	}
	initiatingParticipantPrivateKey := signatureParameters.GetPrivateKeyForIdentity(initiatingParticipant)
	initiatingParticipantSignature := nextInstance.Sign(initiatingParticipantPrivateKey)

	var respondingParticipant *uint
	var respondingParticipantSignature *domain.Signature

	if transition.RespondingParticipant != domain.EmptyParticipantId {
		tmp1 := uint(transition.RespondingParticipant)
		respondingParticipant = &tmp1
		respondingParticipantPrivateKey := signatureParameters.GetPrivateKeyForIdentity(*respondingParticipant)
		tmp2 := nextInstance.Sign(respondingParticipantPrivateKey)
		respondingParticipantSignature = &tmp2
	}

	return State{
		Model:                          model,
		Instance:                       nextInstance,
		Transition:                     transition,
		InitiatingParticipantSignature: initiatingParticipantSignature,
		RespondingParticipantSignature: respondingParticipantSignature,
		InitiatingParticipant:          initiatingParticipant,
		RespondingParticipant:          respondingParticipant,
		ConstraintInput:                constraintInputNotNull,
		InitiatingMessage:              initiatingMessage,
		RespondingMessage:              respondingMessage,
	}
}
