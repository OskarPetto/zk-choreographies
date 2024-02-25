package circuit

import (
	"bytes"
	"execution-service/domain"
	"execution-service/utils"

	"github.com/consensys/gnark-crypto/accumulator/merkletree"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/accumulator/merkle"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type PublicHash struct {
	Hash frontend.Variable `gnark:",public"`
	Salt frontend.Variable
}

type MerkleProof struct {
	MerkleProof merkle.MerkleProof
	Index       frontend.Variable
}

type Authentication struct {
	Signature   eddsa.Signature
	PublicKey   eddsa.PublicKey
	MerkleProof MerkleProof
}

type Instance struct {
	Hash          PublicHash
	Model         frontend.Variable
	TokenCounts   [domain.MaxPlaceCount]frontend.Variable
	PublicKeyRoot frontend.Variable
	MessageHashes [domain.MaxMessageCount]frontend.Variable
}

type ConstraintInput struct {
	Messages [domain.MaxConstraintMessageCount]Message
}

type Message struct {
	IntegerMessage frontend.Variable
	Instance       frontend.Variable
	Salt           frontend.Variable
}

type Constraint struct {
	Coefficients       [domain.MaxConstraintMessageCount]frontend.Variable
	MessageIds         [domain.MaxConstraintMessageCount]frontend.Variable
	Offset             frontend.Variable
	ComparisonOperator frontend.Variable
}

type Transition struct {
	IncomingPlaces        [domain.MaxBranchingFactor]frontend.Variable
	OutgoingPlaces        [domain.MaxBranchingFactor]frontend.Variable
	InitiatingParticipant frontend.Variable
	RespondingParticipant frontend.Variable
	InitiatingMessage     frontend.Variable
	RespondingMessage     frontend.Variable
	Constraint            Constraint
	MerkleProof           MerkleProof
}

type Model struct {
	Salt             frontend.Variable
	PlaceCount       frontend.Variable
	ParticipantCount frontend.Variable
	MessageCount     frontend.Variable
	StartPlaces      [domain.MaxStartPlaceCount]frontend.Variable
	EndPlaceRoot     frontend.Variable
	TransitionRoot   frontend.Variable
}

func ToAuthentication(instance domain.Instance, signature domain.Signature) Authentication {
	var signatureValue eddsa.Signature
	signatureValue.Assign(twistededwards.BN254, signature.Value)

	var buf bytes.Buffer
	for _, publicKey := range instance.PublicKeys {
		hash := publicKey.ComputeHash()
		buf.Write(hash.Value[:])
	}
	for i := len(instance.PublicKeys); i < domain.MaxParticipantCount; i++ {
		publicKey := domain.OutOfBoundsPublicKey()
		hash := publicKey.ComputeHash()
		buf.Write(hash.Value[:])
	}
	participantId := domain.EmptyParticipantId
	for i, publicKey := range instance.PublicKeys {
		if bytes.Equal(publicKey.Value, signature.PublicKey.Value) {
			participantId = domain.ParticipantId(i)
		}
	}
	merkleRoot, proofPath, _, err := merkletree.BuildReaderProof(&buf, hash.MIMC_BN254.New(), fr.Bytes, uint64(participantId))
	utils.PanicOnError(err)
	var merkeProof merkle.MerkleProof
	merkeProof.RootHash = merkleRoot
	merkeProof.Path = make([]frontend.Variable, domain.MaxParticipantDepth+1)
	for i := 0; i < domain.MaxParticipantDepth+1; i++ {
		merkeProof.Path[i] = proofPath[i]
	}
	return Authentication{
		Signature: signatureValue,
		PublicKey: fromPublicKey(signature.PublicKey),
		MerkleProof: MerkleProof{
			MerkleProof: merkeProof,
			Index:       participantId,
		},
	}
}

func FromInstance(instance domain.Instance) Instance {
	var tokenCounts [domain.MaxPlaceCount]frontend.Variable
	for i, tokenCount := range instance.TokenCounts {
		tokenCounts[i] = tokenCount
	}
	for i := len(instance.TokenCounts); i < domain.MaxPlaceCount; i++ {
		tokenCounts[i] = domain.OutOfBoundsTokenCount
	}
	tree := merkletree.New(hash.MIMC_BN254.New())
	for _, publicKey := range instance.PublicKeys {
		hash := publicKey.ComputeHash()
		tree.Push(hash.Value[:])
	}
	for i := len(instance.PublicKeys); i < domain.MaxParticipantCount; i++ {
		publicKey := domain.OutOfBoundsPublicKey()
		hash := publicKey.ComputeHash()
		tree.Push(hash.Value[:])
	}
	var messageHashes [domain.MaxMessageCount]frontend.Variable
	for i, messageHash := range instance.MessageHashes {
		messageHashes[i] = fromHash(messageHash)
	}
	for i := len(instance.MessageHashes); i < domain.MaxMessageCount; i++ {
		messageHashes[i] = fromHash(domain.OutOfBoundsHash())
	}
	return Instance{
		Hash:          toPublicHash(instance.SaltedHash),
		Model:         fromHash(instance.Model),
		TokenCounts:   tokenCounts,
		PublicKeyRoot: tree.Root(),
		MessageHashes: messageHashes,
	}
}

func fromPublicKey(publicKey domain.PublicKey) eddsa.PublicKey {
	var eddsaPublicKey eddsa.PublicKey
	eddsaPublicKey.Assign(twistededwards.BN254, publicKey.Value[:])
	return eddsaPublicKey
}

func FromModel(model domain.Model) Model {
	var startPlaces [domain.MaxStartPlaceCount]frontend.Variable
	for i, startPlace := range model.StartPlaces {
		startPlaces[i] = startPlace
	}
	for i := len(model.StartPlaces); i < domain.MaxStartPlaceCount; i++ {
		startPlaces[i] = domain.OutOfBoundsPlaceId
	}
	endPlaceTree := merkletree.New(hash.MIMC_BN254.New())
	for _, endPlace := range model.EndPlaces {
		bytes := domain.Uint16ToBytes(endPlace)
		endPlaceTree.Push(bytes[:])
	}
	for i := len(model.EndPlaces); i < domain.MaxEndPlaceCount; i++ {
		bytes := domain.Uint16ToBytes(domain.OutOfBoundsPlaceId)
		endPlaceTree.Push(bytes[:])
	}
	transitionTree := merkletree.New(hash.MIMC_BN254.New())
	for _, transition := range model.Transitions {
		hash := transition.ComputeHash()
		transitionTree.Push(hash.Value[:])
	}
	for i := len(model.Transitions); i < domain.MaxTransitionCount; i++ {
		transition := domain.OutOfBoundsTransition()
		hash := transition.ComputeHash()
		transitionTree.Push(hash.Value[:])
	}
	return Model{
		Salt:             fromBytes(model.Hash.Salt),
		PlaceCount:       model.PlaceCount,
		ParticipantCount: model.ParticipantCount,
		MessageCount:     model.MessageCount,
		StartPlaces:      startPlaces,
		EndPlaceRoot:     endPlaceTree.Root(),
		TransitionRoot:   transitionTree.Root(),
	}
}

func ToEndPlaceProof(model domain.Model, instance domain.Instance) MerkleProof {
	var buf bytes.Buffer
	index := 0
	for i, endPlace := range model.EndPlaces {
		if instance.TokenCounts[endPlace] == 1 {
			index = i
		}
		bytes := domain.Uint16ToBytes(endPlace)
		buf.Write(bytes[:])
	}
	for i := len(model.EndPlaces); i < domain.MaxEndPlaceCount; i++ {
		bytes := domain.Uint16ToBytes(domain.OutOfBoundsPlaceId)
		buf.Write(bytes[:])
	}
	merkleRoot, proofPath, _, err := merkletree.BuildReaderProof(&buf, hash.MIMC_BN254.New(), fr.Bytes, uint64(index))
	utils.PanicOnError(err)
	var merkeProof merkle.MerkleProof
	merkeProof.RootHash = merkleRoot
	merkeProof.Path = make([]frontend.Variable, domain.MaxEndPlaceDepth+1)
	for i := 0; i < domain.MaxEndPlaceDepth+1; i++ {
		merkeProof.Path[i] = proofPath[i]
	}
	return MerkleProof{
		MerkleProof: merkeProof,
		Index:       index,
	}
}

func ToTransition(model domain.Model, transition domain.Transition) Transition {
	var incomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i, incomingPlace := range transition.IncomingPlaces {
		incomingPlaces[i] = incomingPlace
	}
	for i := len(transition.IncomingPlaces); i < domain.MaxBranchingFactor; i++ {
		incomingPlaces[i] = domain.OutOfBoundsPlaceId
	}
	var outgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i, outgoingPlace := range transition.OutgoingPlaces {
		outgoingPlaces[i] = outgoingPlace
	}
	for i := len(transition.OutgoingPlaces); i < domain.MaxBranchingFactor; i++ {
		outgoingPlaces[i] = domain.OutOfBoundsPlaceId
	}
	index := domain.MaxTransitionCount
	var buf bytes.Buffer
	for i, modelTransition := range model.Transitions {
		if modelTransition.Id == transition.Id {
			index = i
		}
		hash := modelTransition.ComputeHash()
		buf.Write(hash.Value[:])
	}
	for i := len(model.Transitions); i < domain.MaxTransitionCount; i++ {
		modelTransition := domain.OutOfBoundsTransition()
		hash := modelTransition.ComputeHash()
		buf.Write(hash.Value[:])
	}
	merkleRoot, proofPath, _, err := merkletree.BuildReaderProof(&buf, hash.MIMC_BN254.New(), fr.Bytes, uint64(index))
	utils.PanicOnError(err)
	var merkeProof merkle.MerkleProof
	merkeProof.RootHash = merkleRoot
	merkeProof.Path = make([]frontend.Variable, domain.MaxTransitionDepth+1)
	for i := 0; i < domain.MaxTransitionDepth+1; i++ {
		merkeProof.Path[i] = proofPath[i]
	}
	return Transition{
		IncomingPlaces:        incomingPlaces,
		OutgoingPlaces:        outgoingPlaces,
		InitiatingParticipant: transition.InitiatingParticipant,
		RespondingParticipant: transition.RespondingParticipant,
		InitiatingMessage:     transition.InitiatingMessage,
		RespondingMessage:     transition.RespondingMessage,
		Constraint:            fromConstraint(transition.Constraint),
		MerkleProof: MerkleProof{
			MerkleProof: merkeProof,
			Index:       index,
		},
	}
}

func fromConstraint(constraint domain.Constraint) Constraint {
	var coefficients [domain.MaxConstraintMessageCount]frontend.Variable
	for i, coefficient := range constraint.Coefficients {
		coefficients[i] = coefficient
	}
	for i := len(constraint.Coefficients); i < domain.MaxConstraintMessageCount; i++ {
		coefficients[i] = 0
	}
	var messageIds [domain.MaxConstraintMessageCount]frontend.Variable
	for i, messageId := range constraint.MessageIds {
		messageIds[i] = messageId
	}
	for i := len(constraint.MessageIds); i < domain.MaxConstraintMessageCount; i++ {
		messageIds[i] = domain.EmptyMessageId
	}
	return Constraint{
		Coefficients:       coefficients,
		MessageIds:         messageIds,
		Offset:             constraint.Offset,
		ComparisonOperator: constraint.ComparisonOperator,
	}
}

func FromConstraintInput(constraintInput domain.ConstraintInput) ConstraintInput {
	var messages [domain.MaxConstraintMessageCount]Message
	for i, message := range constraintInput.Messages {
		messages[i] = fromMessage(message)
	}
	for i := len(constraintInput.Messages); i < domain.MaxConstraintMessageCount; i++ {
		messages[i] = fromMessage(domain.EmptyMessage())
	}
	return ConstraintInput{
		Messages: messages,
	}
}

func fromMessage(message domain.Message) Message {
	return Message{
		IntegerMessage: message.IntegerMessage,
		Instance:       fromBytes(message.Instance.Value),
		Salt:           fromBytes(message.Hash.Salt),
	}
}

func toPublicHash(hash domain.SaltedHash) PublicHash {
	return PublicHash{
		Hash: fromHash(hash.Hash),
		Salt: fromBytes(hash.Salt),
	}
}

func fromHash(hash domain.Hash) frontend.Variable {
	return fromBytes(hash.Value)
}

func fromBytes(data [fr.Bytes]byte) frontend.Variable {
	fieldElement, err := fr.BigEndian.Element(&data)
	utils.PanicOnError(err)
	return fieldElement
}
