package blockchain

import "time"

type Block struct {
	Index     int64     `json:"index"`
	Hash      string    `json:"hash"`
	PrevHash  string    `json:"prev_hash"`
	Timestamp time.Time `json:"timestamp"`
	Data      []byte    `json:"data"` // TODO data as byte or string?
}
