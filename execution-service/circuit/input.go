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
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/signature/eddsa"
)

type Hash struct {
	Value frontend.Variable `gnark:",public"`
	Salt  frontend.Variable
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
	Hash          Hash
	TokenCounts   [domain.MaxPlaceCount]frontend.Variable
	PublicKeyRoot frontend.Variable
	MessageHashes [domain.MaxMessageCount]frontend.Variable
}

type ConstraintInput struct {
	IntegerMessages [domain.MaxConstraintMessageCount]frontend.Variable
	Salts           [domain.MaxConstraintMessageCount]frontend.Variable
}

type Constraint struct {
	Coefficients       [domain.MaxConstraintMessageCount]frontend.Variable
	MessageIds         [domain.MaxConstraintMessageCount]frontend.Variable
	Offset             frontend.Variable
	ComparisonOperator frontend.Variable
}

type Transition struct {
	IncomingPlaces [domain.MaxBranchingFactor]frontend.Variable
	OutgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	Participant    frontend.Variable
	Message        frontend.Variable
	Constraint     Constraint
	MerkleProof    MerkleProof
}

type Model struct {
	Hash             Hash
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
		buf.Write(hash[:])
	}
	merkleRoot, proofPath, _, err := merkletree.BuildReaderProof(&buf, hash.MIMC_BN254.New(), fr.Bytes, uint64(signature.Participant))
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
			Index:       signature.Participant,
		},
	}
}

func FromInstance(instance domain.Instance) Instance {
	var tokenCounts [domain.MaxPlaceCount]frontend.Variable
	for i, tokenCount := range instance.TokenCounts {
		tokenCounts[i] = tokenCount
	}
	tree := merkletree.New(hash.MIMC_BN254.New())
	for _, publicKey := range instance.PublicKeys {
		hash := publicKey.ComputeHash()
		tree.Push(hash[:])
	}
	var messageHashes [domain.MaxMessageCount]frontend.Variable
	for i, messageHash := range instance.MessageHashes {
		messageHashes[i] = fromBytes(messageHash)
	}
	return Instance{
		Hash:          fromHash(instance.Hash),
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
	endPlaceTree := merkletree.New(hash.MIMC_BN254.New())
	for _, endPlace := range model.EndPlaces {
		bytes := domain.Uint8ToBytes(endPlace)
		endPlaceTree.Push(bytes[:])
	}
	transitionTree := merkletree.New(hash.MIMC_BN254.New())
	for _, transition := range model.Transitions {
		hash := transition.ComputeHash()
		transitionTree.Push(hash[:])
	}
	return Model{
		Hash:             fromHash(model.Hash),
		PlaceCount:       model.PlaceCount,
		ParticipantCount: model.ParticipantCount,
		MessageCount:     model.MessageCount,
		StartPlaces:      startPlaces,
		EndPlaceRoot:     endPlaceTree.Root(),
		TransitionRoot:   transitionTree.Root(),
	}
}

func ToEndPlace(model domain.Model, place domain.PlaceId) MerkleProof {
	var buf bytes.Buffer
	index := domain.MaxEndPlaceCount
	for i, endPlace := range model.EndPlaces {
		if endPlace == place {
			index = i
		}
		bytes := domain.Uint8ToBytes(endPlace)
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
	var outgoingPlaces [domain.MaxBranchingFactor]frontend.Variable
	for i, outgoingPlace := range transition.OutgoingPlaces {
		outgoingPlaces[i] = outgoingPlace
	}
	index := domain.MaxTransitionCount
	var buf bytes.Buffer
	for i, modelTransition := range model.Transitions {
		if modelTransition.Id == transition.Id {
			index = i
		}
		hash := modelTransition.ComputeHash()
		buf.Write(hash[:])
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
		IncomingPlaces: incomingPlaces,
		OutgoingPlaces: outgoingPlaces,
		Participant:    transition.Participant,
		Message:        transition.Message,
		Constraint:     fromConstraint(transition.Constraint),
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
	var messageIds [domain.MaxConstraintMessageCount]frontend.Variable
	for i, messageId := range constraint.MessageIds {
		messageIds[i] = messageId
	}
	return Constraint{
		Coefficients:       coefficients,
		MessageIds:         messageIds,
		Offset:             constraint.Offset,
		ComparisonOperator: constraint.ComparisonOperator,
	}
}

func FromConstraintInput(input domain.ConstraintInput) ConstraintInput {
	var integerMessages [domain.MaxConstraintMessageCount]frontend.Variable
	var salts [domain.MaxConstraintMessageCount]frontend.Variable
	for i, message := range input.Messages {
		integerMessages[i] = message.IntegerMessage
		salts[i] = fromBytes(message.Hash.Salt)
	}
	return ConstraintInput{
		IntegerMessages: integerMessages,
		Salts:           salts,
	}
}

func fromHash(hash domain.Hash) Hash {
	return Hash{
		Value: fromBytes(hash.Value),
		Salt:  fromBytes(hash.Salt),
	}
}

func fromBytes(data [fr.Bytes]byte) frontend.Variable {
	fieldElement, err := fr.BigEndian.Element(&data)
	utils.PanicOnError(err)
	return fieldElement
}

func (merkleProof *MerkleProof) CheckRootHash(api frontend.API, hash frontend.Variable) {
	api.AssertIsEqual(merkleProof.MerkleProof.RootHash, hash)
}

func (merkleProof *MerkleProof) VerifyProof(api frontend.API, m mimc.MiMC) {
	merkleProof.MerkleProof.VerifyProof(api, &m, merkleProof.Index)
}
