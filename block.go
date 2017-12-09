package blockchain

import (
	"time"
)

const GenesisBlockIndex = 0

type Block struct {
	hasher    Hasher
	Index     int64     `json:"index"`
	Hash      string    `json:"hash"`
	PrevHash  string    `json:"prev_hash"`
	Timestamp time.Time `json:"timestamp"`
	Nonce     int       `json:"nonce"`
	Data      []byte    `json:"data"`
}

type Hasher interface {
	GenerateHash(index int64, prevHash string, timestamp time.Time, data []byte) string
}

func GenesisBlock(h Hasher, data []byte) *Block {
	// TODO add proof of work
	t := time.Now()
	return &Block{
		hasher:    h,
		Index:     GenesisBlockIndex,
		Hash:      h.GenerateHash(GenesisBlockIndex, "", t, data),
		PrevHash:  "",
		Timestamp: t,
		Data:      data,
	}
}

func (b *Block) IsGenesisBlock() bool {
	return b.Index == GenesisBlockIndex && b.PrevHash == ""
}

func Mine(tip *Block, data []byte) *Block {
	index := tip.Index + 1
	t := time.Now()
	return &Block{
		hasher:    tip.hasher,
		Index:     index,
		Hash:      tip.hasher.GenerateHash(index, tip.Hash, t, data),
		PrevHash:  tip.Hash,
		Timestamp: t,
		Data:      data,
	}
}
