package blockchain

import (
	"testing"
	"time"

	"github.com/siklol/blockchain/hasher"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulGenesisBlock(t *testing.T) {
	// given
	gt := time.Date(2017, 12, 9, 12, 0, 0, 0, &time.Location{})

	// when
	g := GenesisBlock(hasher.Sha256Hasher, gt, []byte("genesis block!"))

	// then
	assert.NotEmpty(t, g)
	assert.Equal(t, g.Hash, "LE2Z/RqRua8vNB7yS2tYRuoNes5dCuiYmywtKQpVoBo=")
}
