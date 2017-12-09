package hasher

import (
	"crypto"
	"encoding/base64"
	"strconv"
	"time"
)

var (
	Sha256Hasher = Hasher{}
)

type Hasher struct {
}

func (ha Hasher) GenerateHash(index int, prevHash string, timestamp time.Time, data []byte) string {
	h := crypto.SHA256.New()
	h.Write([]byte(strconv.Itoa(index) + prevHash + string(timestamp.Unix()))) // TODO check for errors
	h.Write(data)
	hashByte := h.Sum(nil)

	return base64.StdEncoding.EncodeToString(hashByte)
}
