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
	ConditionInput                 domain.ConditionInput
	InitiatingMessage              *domain.Message
	RespondingMessage              *domain.Message
}

func GetModel2States(signatureParameters parameters.SignatureParameters) []State {
	model2 := GetModel2()

	state0 := instantiateModel(signatureParameters, model2)
	state1 := executeTransition(signatureParameters, model2, 0, state0.Instance, nil, nil, nil)
	order, err := domain.NewInitiatingIntegerMessage(state1.Instance, model2.Transitions[2], 2)
	utils.PanicOnError(err)
	stock, err := domain.NewRespondingIntegerMessage(state1.Instance, model2.Transitions[2], 20)
	utils.PanicOnError(err)
	state2 := executeTransition(signatureParameters, model2, 2, state1.Instance, nil, &order, &stock)
	confirm, err := domain.NewInitiatingBytesMessage(state2.Instance, model2.Transitions[7], []byte("confirm"))
	utils.PanicOnError(err)
	conditionInput1 := domain.ConditionInput{
		Messages: []domain.Message{
			order, stock,
		},
	}
	state3 := executeTransition(signatureParameters, model2, 7, state2.Instance, &conditionInput1, &confirm, nil)
	invoice, err := domain.NewInitiatingBytesMessage(state3.Instance, model2.Transitions[6], []byte("invoice"))
	utils.PanicOnError(err)
	payment, err := domain.NewRespondingBytesMessage(state3.Instance, model2.Transitions[6], []byte("payment"))
	utils.PanicOnError(err)
	state4 := executeTransition(signatureParameters, model2, 6, state3.Instance, nil, &invoice, &payment)
	shippingAddress, err := domain.NewInitiatingBytesMessage(state4.Instance, model2.Transitions[5], []byte("shipping_address"))
	utils.PanicOnError(err)
	product, err := domain.NewInitiatingBytesMessage(state4.Instance, model2.Transitions[5], []byte("product"))
	utils.PanicOnError(err)
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

	initiatingParticipantPrivateKey, err := signatureParameters.GetPrivateKeyForIdentity(0)
	utils.PanicOnError(err)
	initiatingParticipantSignature := instance.Sign(initiatingParticipantPrivateKey)
	return State{
		Model:                          model,
		Instance:                       instance,
		Transition:                     domain.OutOfBoundsTransition(),
		InitiatingParticipantSignature: initiatingParticipantSignature,
	}
}

func executeTransition(signatureParameters parameters.SignatureParameters, model domain.Model, transitionIndex uint, currentInstance domain.Instance, conditionInput *domain.ConditionInput, initiatingMessage *domain.Message, respondingMessage *domain.Message) State {

	var transition domain.Transition
	if transitionIndex < domain.MaxTransitionCount {
		transition = model.Transitions[transitionIndex]
	} else {
		transition = domain.OutOfBoundsTransition()
	}
	conditionInputNotNull := domain.EmptyConditionInput()
	if conditionInput != nil {
		conditionInputNotNull = *conditionInput
	}

	nextInstance, err := currentInstance.ExecuteTransition(transition, conditionInputNotNull, initiatingMessage, respondingMessage)
	utils.PanicOnError(err)

	initiatingParticipant := uint(0)
	if transition.InitiatingParticipant != domain.EmptyParticipantId {
		initiatingParticipant = uint(transition.InitiatingParticipant)
	}
	initiatingParticipantPrivateKey, err := signatureParameters.GetPrivateKeyForIdentity(initiatingParticipant)
	utils.PanicOnError(err)
	initiatingParticipantSignature := nextInstance.Sign(initiatingParticipantPrivateKey)

	var respondingParticipant *uint
	var respondingParticipantSignature *domain.Signature

	if transition.RespondingParticipant != domain.EmptyParticipantId {
		tmp1 := uint(transition.RespondingParticipant)
		respondingParticipant = &tmp1
		respondingParticipantPrivateKey, err := signatureParameters.GetPrivateKeyForIdentity(*respondingParticipant)
		utils.PanicOnError(err)
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
		ConditionInput:                 conditionInputNotNull,
		InitiatingMessage:              initiatingMessage,
		RespondingMessage:              respondingMessage,
	}
}
