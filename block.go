package blockchain

import (
	"errors"
	"time"
)

const GenesisBlockIndex = 0
const Salt = "64kjwsfhgm2w46ktwe6tulkdgfa345werzh1q435jhwrtzk5e37lk" // FIXME

var (
	ErrInvalidBlockHash          = errors.New("invalid block hash")
	ErrInvalidBlock              = errors.New("invalid block")
	ErrMissingPreviousBlock      = errors.New("missing previous block")
	ErrCannotDestroyGenesisBlock = errors.New("destroying from genesis block is forbidden")
	ErrInvalidDestroyBlock       = errors.New("cannot destroy block at index")
)

type Block struct {
	Index     int64     `json:"index"`
	Hash      string    `json:"hash"`
	PrevHash  string    `json:"prev_hash"`
	Timestamp time.Time `json:"timestamp"`
	Nonce     int64     `json:"nonce"`
	Data      []byte    `json:"data"`
}

type Hasher interface {
	GenerateHash(index int64, nonce int64, prevHash string, timestamp time.Time, data []byte) string
}

type ProofOfWork interface {
	Proof(data []byte, t time.Time, salt string) int64
	Verify(data []byte, nonce int64, t time.Time, salt string) bool
}

func GenesisBlock(h Hasher, p ProofOfWork, data []byte, t time.Time) *Block {
	nonce := p.Proof(data, t, Salt)

	return &Block{
		Index:     GenesisBlockIndex,
		Nonce:     nonce,
		Hash:      h.GenerateHash(GenesisBlockIndex, nonce, "", t, data),
		PrevHash:  "",
		Timestamp: t,
		Data:      data,
	}
}

func (b *Block) IsGenesisBlock() bool {
	return b.Index == GenesisBlockIndex && b.PrevHash == ""
}

func Mine(tip *Block, data []byte, hasher Hasher, proof ProofOfWork) (*Block, error) {
	index := tip.Index + 1
	t := time.Now()

	if !tip.IsValidHash(hasher, tip.Hash) {
		return nil, ErrInvalidBlockHash
	}

	nonce := proof.Proof(data, t, Salt)

	b := &Block{
		Index:     index,
		Nonce:     nonce,
		Hash:      hasher.GenerateHash(index, nonce, tip.Hash, t, data),
		PrevHash:  tip.Hash,
		Timestamp: t,
		Data:      data,
	}

	return b, nil
}

func (b *Block) IsValidHash(hasher Hasher, hash string) bool {
	challenge := hasher.GenerateHash(b.Index, b.Nonce, b.PrevHash, b.Timestamp, b.Data)
	return hash == challenge
}

func (b *Block) VerifyProofOfWork(proof ProofOfWork, nonce int64) bool {
	return proof.Verify(b.Data, nonce, b.Timestamp, Salt)
}

func (b *Block) Verify(hasher Hasher, proof ProofOfWork) bool {
	return b.VerifyProofOfWork(proof, b.Nonce) && b.IsValidHash(hasher, b.Hash)
}
