package blockchain

import (
	"time"

	"github.com/siklol/blockchain/hasher"
	"github.com/siklol/blockchain/proofofwork"
)

type HashType string
type ProofOfWorkType string

const Sha256 HashType = "sha256"
const Hashcash = "hashcash"

type Blockchain struct {
	blocks      map[string]*Block
	indexToHash map[int64]string
	tip         *Block
	hasher      Hasher
	proof       ProofOfWork
}

func NewBlockchain(ht HashType, pow ProofOfWorkType, genesisMsg []byte, genesisTimestamp time.Time) *Blockchain {
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
	b := GenesisBlock(hashingAlgo, powAlgo, genesisMsg, genesisTimestamp)
	indexMap := make(map[int64]string)

	blocks[b.Hash] = b
	indexMap[0] = b.Hash

	return &Blockchain{
		blocks:      blocks,
		tip:         b,
		indexToHash: indexMap,
		hasher:      hashingAlgo,
		proof:       powAlgo,
	}
}

func (c *Blockchain) Tip() *Block {
	return c.tip
}

func (c *Blockchain) Mine(d []byte) error {
	newTip, err := Mine(c.tip, d, c.hasher, c.proof)

	if err != nil {
		return err
	}

	c.tip = newTip
	c.blocks[newTip.Hash] = newTip
	c.indexToHash[newTip.Index] = newTip.Hash
	return nil
}

func (c *Blockchain) PreviousBlock(b *Block) *Block {
	if b == nil {
		return nil
	}

	prev, ok := c.blocks[b.PrevHash]

	if !ok {
		return nil
	}

	return prev
}
func (c *Blockchain) Blocks() []*Block {
	blocks := []*Block{}

	tip := c.Tip()
	for {
		prevBlock := c.PreviousBlock(tip)

		if tip == nil {
			break
		}

		blocks = append(blocks, tip)
		tip = prevBlock
	}

	return blocks
}
func (c *Blockchain) Genesis() *Block {
	return c.blocks[c.indexToHash[0]]
}
func (c *Blockchain) Exists(block *Block) bool {
	_, ok := c.blocks[block.Hash]

	return ok
}
func (c *Blockchain) BlockAtIndex(i int64) *Block {
	hash, ok := c.indexToHash[i]

	if ok {
		return c.blocks[hash]
	}

	return nil
}
func (c *Blockchain) Append(block *Block) error {
	if !block.Verify(c.hasher, c.proof) {
		return ErrInvalidBlock
	}

	if c.tip.Hash != block.PrevHash {
		return ErrMissingPreviousBlock
	}

	c.tip = block
	c.blocks[block.Hash] = block
	c.indexToHash[block.Index] = block.Hash
	return nil
}
