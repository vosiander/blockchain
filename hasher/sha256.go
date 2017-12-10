package hasher

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/siklol/blockchain/conversion"
)

var (
	Sha256 = hasher{}
)

type hasher struct {
}

func (ha hasher) GenerateHash(index int64, nonce int64, prevHash string, timestamp time.Time, data []byte) string {
	b := bytes.Join(
		[][]byte{
			conversion.IntToHex(index),
			conversion.IntToHex(nonce),
			[]byte(prevHash),
			conversion.IntToHex(timestamp.Unix()),
		},
		data,
	)

	hashByte := sha256.Sum256(b)

	return base64.StdEncoding.EncodeToString(hashByte[:])
}
