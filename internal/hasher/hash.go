package hasher

import "crypto/sha256"

func Hash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return string(h.Sum(nil))
}
