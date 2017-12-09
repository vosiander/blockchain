package blockchain

import (
	"testing"

	"github.com/siklol/blockchain/hasher"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulGenesisBlock(t *testing.T) {
	// given
	data := []byte("genesis block!")

	// when
	g := GenesisBlock(hasher.Sha256, data)

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
	g := GenesisBlock(hasher.Sha256, data)
	c := Mine(g, secondBlockData)

	// then
	assert.NotEmpty(t, g)
	assert.NotEmpty(t, c)
	assert.NotEmpty(t, g.Hash)
	assert.Equal(t, g.Hash, c.PrevHash)
	assert.NotEmpty(t, c)
	assert.True(t, g.IsGenesisBlock())
	assert.False(t, c.IsGenesisBlock())
}
