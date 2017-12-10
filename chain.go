package blockchain

import (
	"github.com/siklol/blockchain/hasher"
	"github.com/siklol/blockchain/proofofwork"
)

type HashType string

const Sha256 HashType = "sha256"

type ProofOfWorkType string

const Hashcash = "hashcash"

type Blockchain struct {
	blocks      map[string]*Block
	tip         *Block
	hashingAlgo HashType

	pow ProofOfWorkType
}

func NewBlockchain(ht HashType, pow ProofOfWorkType, genesisMsg []byte) *Blockchain {
	var hashingAlgo Hasher
	switch ht {
	case Sha256:
		hashingAlgo = hasher.Sha256
	default:
		hashingAlgo = hasher.Sha256
	}

	var powAlgo ProofOfWork
	switch pow {
	case Hashcash:
		powAlgo = proofofwork.HashCash
	default:
		powAlgo = proofofwork.HashCash
	}

	blocks := make(map[string]*Block)
	b := GenesisBlock(hashingAlgo, powAlgo, genesisMsg)

	blocks[b.Hash] = b

	return &Blockchain{
		blocks:      blocks,
		tip:         b,
		hashingAlgo: ht,
		pow:         pow,
	}
}

func (c *Blockchain) Tip() *Block {
	return c.tip
}

func (c *Blockchain) Mine(d []byte) error {
	newTip, err := Mine(c.tip, d)

	if err != nil {
		return err
	}

	c.tip = newTip
	c.blocks[newTip.Hash] = newTip
	return nil
}

func (c *Blockchain) PreviousBlock(b *Block) *Block {
	prev, ok := c.blocks[b.PrevHash]

	if !ok {
		return nil
	}

	return prev
}
