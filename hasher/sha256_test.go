package hasher_test

import (
	"testing"
	"time"

	"github.com/siklol/blockchain/hasher"
	"github.com/stretchr/testify/assert"
)

type hasherTestStruct struct {
	index     int64
	prevHash  string
	timestamp time.Time
	data      []byte
}

func TestHashGeneration(t *testing.T) {
	// given
	tx := time.Now()
	data := []struct {
		first     hasherTestStruct
		challenge hasherTestStruct
	}{
		{
			first:     hasherTestStruct{0, "", tx, []byte("This is a test")},
			challenge: hasherTestStruct{1, "", tx, []byte("This is a test")},
		},
		{
			first:     hasherTestStruct{0, "", tx, []byte("This is a test")},
			challenge: hasherTestStruct{0, "abc", tx, []byte("This is a test")},
		},
		{
			first:     hasherTestStruct{0, "", tx, []byte("This is a test")},
			challenge: hasherTestStruct{0, "", tx.Add(-1 * 5 * time.Minute), []byte("This is a test")},
		},
		{
			first:     hasherTestStruct{0, "", tx, []byte("This is a test")},
			challenge: hasherTestStruct{0, "", tx, []byte("That is really one")},
		},
	}

	for _, d := range data {
		// when
		actual := hasher.Sha256.GenerateHash(d.first.index, d.first.prevHash, d.first.timestamp, d.first.data)
		challange := hasher.Sha256.GenerateHash(d.challenge.index, d.challenge.prevHash, d.challenge.timestamp, d.challenge.data)

		// then
		assert.NotEqual(t, "", actual)
		assert.NotEqual(t, actual, challange)
	}
}
