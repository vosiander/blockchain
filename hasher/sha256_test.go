package hasher_test

import (
	"testing"
	"time"

	"github.com/siklol/blockchain/hasher"
	"github.com/stretchr/testify/assert"
)

type hasherTestStruct struct {
	index     int64
	nonce     int64
	prevHash  string
	timestamp time.Time
	data      []byte
}

func TestHashGeneration(t *testing.T) {
	// given
	tx := time.Now()
	data := []struct {
		a hasherTestStruct
		b hasherTestStruct
	}{
		{
			a: hasherTestStruct{0, 1, "", tx, []byte("This is a test")},
			b: hasherTestStruct{1, 1, "", tx, []byte("This is a test")},
		},
		{
			a: hasherTestStruct{0, 1, "", tx, []byte("This is a test")},
			b: hasherTestStruct{0, 1, "abc", tx, []byte("This is a test")},
		},
		{
			a: hasherTestStruct{0, 1, "", tx, []byte("This is a test")},
			b: hasherTestStruct{0, 1, "", tx.Add(-1 * 5 * time.Minute), []byte("This is a test")},
		},
		{
			a: hasherTestStruct{0, 1, "", tx, []byte("This is a test")},
			b: hasherTestStruct{0, 1, "", tx, []byte("That is really one")},
		},
	}

	for _, d := range data {
		// when
		actual := hasher.Sha256.GenerateHash(d.a.index, d.a.nonce, d.a.prevHash, d.a.timestamp, d.a.data)
		challange := hasher.Sha256.GenerateHash(d.b.index, d.b.nonce, d.b.prevHash, d.b.timestamp, d.b.data)

		// then
		assert.NotEqual(t, "", actual)
		assert.NotEqual(t, actual, challange)
	}
}
