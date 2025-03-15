package hasher

import "crypto/sha256"

func Hash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	return string(h.Sum(nil))
}
