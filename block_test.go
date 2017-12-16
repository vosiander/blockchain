package blockchain

import (
	"testing"

	"time"

	"github.com/siklol/blockchain/hasher"
	"github.com/siklol/blockchain/proofofwork"
	"github.com/stretchr/testify/assert"
)

var (
	proofOfWorkForTest = proofofwork.HashCash
	genesisTimestamp   = time.Date(2017, 12, 10, 12, 0, 0, 0, time.UTC)
)

func init() {
	proofOfWorkForTest.Difficulty = 16 // use a simple difficulty for fast test results
}

func TestSuccessfulGenesisBlock(t *testing.T) {
	// given
	data := []byte("genesis block!")

	// when
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data, genesisTimestamp)

	// then
	assert.NotEmpty(t, g)
	assert.NotEmpty(t, g.Hash)
	assert.True(t, g.IsGenesisBlock())
}

func TestSuccessfulMineBlock(t *testing.T) {
	// given
	data := []byte("genesis block!")
	secondBlockData := []byte("blockchain rocks!")

	// when
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data, genesisTimestamp)
	c, err := Mine(g, secondBlockData, hasher.Sha256, proofofwork.HashCash)

	// then
	assert.NoError(t, err)
	if err != nil {
		t.FailNow()
	}

	assert.NotEmpty(t, g)
	assert.NotEmpty(t, c)
	assert.NotEmpty(t, g.Hash)
	assert.Equal(t, g.Hash, c.PrevHash)
	assert.NotEmpty(t, c)
	assert.True(t, g.IsGenesisBlock())
	assert.False(t, c.IsGenesisBlock())
}

func TestInvalidHash(t *testing.T) {
	// given
	data := []byte("genesis block!")
	secondBlockData := []byte("blockchain rocks!")

	// when
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data, genesisTimestamp)
	g.Timestamp = g.Timestamp.Add(-1 * 20 * time.Hour)
	_, err := Mine(g, secondBlockData, hasher.Sha256, proofofwork.HashCash)

	// then
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidBlockHash, err)
}

func TestProofOfWork(t *testing.T) {
	// given
	data := []byte("genesis block!")
	secondBlockData := []byte("blockchain rocks!")

	// when
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data, genesisTimestamp)
	c, _ := Mine(g, secondBlockData, hasher.Sha256, proofofwork.HashCash)

	// then
	assert.True(t, c.VerifyProofOfWork(proofofwork.HashCash, c.Nonce))
}

func TestIdenticalGenesisBlocks(t *testing.T) {
	// given
	data := []byte("genesis block!")

	// when
	firstG := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data, genesisTimestamp)
	secondG := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data, genesisTimestamp)

	// then
	assert.Equal(t, firstG.Hash, secondG.Hash)
}
