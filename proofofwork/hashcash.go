package proofofwork

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
	"time"

	"github.com/siklol/blockchain/conversion"
)

const defaultDifficulty = 18
const maxNonce = math.MaxInt64

func init() {
	HashCash = hashCash{
		Difficulty: defaultDifficulty,
	}
}

var (
	HashCash hashCash
)

type hashCash struct {
	Difficulty int
	exec       string
}

func (h hashCash) Proof(data []byte, t time.Time, salt string) int64 {
	var hashInt big.Int
	var hash [32]byte
	nonce := int64(0)

	target := h.difficulty()

	for nonce < maxNonce {
		concatedData := h.concatData(data, []byte(salt), t, nonce)
		hash = sha256.Sum256(concatedData)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce
}

func (h hashCash) Verify(data []byte, nonce int64, t time.Time, salt string) bool {
	var hashInt big.Int
	target := h.difficulty()

	conctedData := h.concatData(data, []byte(salt), t, nonce)
	hash := sha256.Sum256(conctedData)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(target) == -1

	return isValid
}

func (h hashCash) difficulty() *big.Int {
	target := big.NewInt(1)
	return target.Lsh(target, uint(256-h.Difficulty))
}

func (h hashCash) concatData(data []byte, salt []byte, t time.Time, nonce int64) []byte {
	return bytes.Join(
		[][]byte{
			data,
			salt,
			conversion.IntToHex(t.Unix()),
			conversion.IntToHex(int64(h.Difficulty)),
			conversion.IntToHex(nonce),
		},
		[]byte{},
	)
}
