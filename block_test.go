package blockchain

import (
	"testing"

	"time"

	"github.com/siklol/blockchain/hasher"
	"github.com/siklol/blockchain/proofofwork"
	"github.com/stretchr/testify/assert"
)

var proofOfWorkForTest = proofofwork.HashCash

func init() {
	proofOfWorkForTest.Difficulty = 16 // use a simple difficulty for fast test results
}

func TestSuccessfulGenesisBlock(t *testing.T) {
	// given
	data := []byte("genesis block!")

	// when
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data)

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
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data)
	c, err := Mine(g, secondBlockData)

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
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data)
	g.Timestamp = g.Timestamp.Add(-1 * 20 * time.Hour)
	_, err := Mine(g, secondBlockData)

	// then
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidBlockHash, err)
}

func TestProofOfWork(t *testing.T) {
	// given
	data := []byte("genesis block!")
	secondBlockData := []byte("blockchain rocks!")

	// when
	g := GenesisBlock(hasher.Sha256, proofOfWorkForTest, data)
	c, _ := Mine(g, secondBlockData)

	// then
	assert.True(t, c.VerifyProofOfWork(c.Nonce))
}
