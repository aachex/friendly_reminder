package hasher

import "crypto/sha256"

func Hash(data []byte, hashed *[]byte) {
	h := sha256.New()
	h.Write(data)
	*hashed = h.Sum(nil)
}
