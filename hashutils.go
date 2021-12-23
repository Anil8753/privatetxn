package privatetxn

import "crypto/sha256"

func getHash(input []byte) []byte {
	hash := sha256.New()
	hash.Write(input)
	return hash.Sum(nil)
}
