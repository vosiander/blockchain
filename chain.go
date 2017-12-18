package blockchain

import (
	"log"
	"time"

	"sync"

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
	mu          *sync.Mutex
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
		mu:          &sync.Mutex{},
	}
}

func (c *Blockchain) Tip() *Block {
	return c.tip
}

func (c *Blockchain) Mine(d []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

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
	c.mu.Lock()
	defer c.mu.Unlock()

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

func (c *Blockchain) DestroyBlocksFromIndex(index int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index == 0 {
		return ErrCannotDestroyGenesisBlock
	}

	_, ok := c.indexToHash[index]
	if !ok {
		return ErrInvalidDestroyBlock
	}

	newTipHash := c.blocks[c.indexToHash[index]].PrevHash
	c.tip = c.blocks[newTipHash]

	for true {
		hash, ok := c.indexToHash[index]
		if !ok {
			break
		}

		delete(c.blocks, hash)
		delete(c.indexToHash, index)

		index++
	}

	return nil
}

func (c *Blockchain) ReplaceChainFromIndex(index int64, blocks []*Block) error {

	if index <= c.tip.Index {
		if err := c.DestroyBlocksFromIndex(index); err != nil {
			return err
		}
	}

	for _, b := range blocks {
		if err := c.Append(b); err != nil {
			log.Println("Could not append block to current blockchain. error: " + err.Error())
			return nil
		}
	}

	return nil
}
