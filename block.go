package blockchain

import (
	"time"
)

const GenesisBlockIndex = 0

type Block struct {
	Index     int64     `json:"index"`
	Hash      string    `json:"hash"`
	PrevHash  string    `json:"prev_hash"`
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"` // TODO data as byte or string?
}

type Hasher interface {
	GenerateHash(index int, prevHash string, timestamp time.Time, data []byte) string
}

func GenesisBlock(h Hasher, t time.Time, data []byte) *Block {
	return &Block{
		Index:     GenesisBlockIndex,
		Hash:      h.GenerateHash(GenesisBlockIndex, "", t, data),
		PrevHash:  "",
		Timestamp: t,
		Data:      data,
	}
}
