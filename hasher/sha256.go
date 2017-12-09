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

func (ha hasher) GenerateHash(index int64, prevHash string, t time.Time, data []byte) string {
	b := bytes.Join(
		[][]byte{
			conversion.IntToHex(index),
			[]byte(prevHash),
			conversion.IntToHex(t.Unix()),
		},
		data,
	)

	hashByte := sha256.Sum256(b)

	return base64.StdEncoding.EncodeToString(hashByte[:])
}
